//
//  ContentView.swift
//  KMusic
//
//  Created by Kevin Peng on 2026-02-02.
//

import SwiftUI
import OpenAPIRuntime
import OpenAPIURLSession

struct ContentView: View {
    @State private var id = ""
    @State private var song = Song.empty
    @State private var status = ""

    let client: Client

    init() {
        self.client = Client(
            serverURL: try! Servers.Server1.url(),
            configuration: .init(dateTranscoder: .iso8601WithFractionalSeconds),
            transport: URLSessionTransport(),
            middlewares: [
                APIKeyMiddleware(apiKey: "dhowcWZ3YMdcg2KueEtKG6fCeXkMeRsN")
            ]
        )
    }

    func getSong(id: String) async throws {
        let response = try await client.SongsShow(Operations.SongsShow.Input(path: .init(id: id)))

        switch response {
        case let .ok(body):
            guard let song = try body.body.json.song else {
                status = "Not structured correctly"
                return
            }
            guard let url = URL(string: song.audioUrl) else {
                status = "Song URL malformed"
                return
            }
            self.song = Song(song, url: url)
        case .notFound:
            status = "Not Found"
        case .internalServerError:
            status = "Internal Server Error"
        case .undocumented:
            status = "Undocumented"
        case .unauthorized(_):
            status = "Unauthorized"
        }
    }

    var body: some View {
        VStack {
            Text("Status: \(status)")
            TextField("ID", text: $id)
            Text("Song title: \(song.title)")
            Text("Created at: \(song.createdAt, format: .dateTime)")
            Text("Updated at: \(song.updatedAt, format: .dateTime)")
            Text("Audio URL: \(song.audioURL.absoluteString)")
            Button("Look Up") {
                Task { try? await getSong(id: id) }
            }
        }
        .padding()
    }
}

#Preview {
    ContentView()
}
