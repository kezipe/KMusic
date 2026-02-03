//
//  Song.swift
//  KMusic
//
//  Created by Kevin Peng on 2026-02-03.
//
import Foundation

struct Song {
    let id: Int
    let title: String
    let createdAt: Date
    let updatedAt: Date
    let audioURL: URL

}

extension Song {
    init(_ dto: Operations.SongsShow.Output.Ok.Body.jsonPayload.songPayload, url: URL) {
        self.id = dto.id
        self.title = dto.title
        self.createdAt = dto.createdAt
        self.updatedAt = dto.updatedAt
        self.audioURL = url
    }

    static let empty = Song(id: 0, title: "", createdAt: .now, updatedAt: .now, audioURL: URL(string: "www.ppnn.ca")!)
}
