//
//  APIKeyMiddleware.swift
//  KMusic
//
//  Created by Kevin Peng on 2026-02-02.
//

import OpenAPIRuntime
import HTTPTypes
import Foundation

struct APIKeyMiddleware: ClientMiddleware {

    let apiKey: String

    func intercept(
        _ request: HTTPRequest,
        body: HTTPBody?,
        baseURL: URL,
        operationID: String,
        next: @Sendable (
            HTTPRequest,
            HTTPBody?,
            URL
        ) async throws -> (HTTPResponse, HTTPBody?)
    ) async throws -> (HTTPResponse, HTTPBody?) {

        var request = request
        request.headerFields[HTTPField.Name("Api-Key")!] = apiKey

        return try await next(request, body, baseURL)
    }
}

