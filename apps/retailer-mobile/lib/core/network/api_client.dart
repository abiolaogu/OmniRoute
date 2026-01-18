/// OmniRoute Ecosystem - Network Layer
/// Comprehensive API client with authentication, caching, and error handling

import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:omniroute_ecosystem/core/constants/app_constants.dart';

// ============================================================================
// API CLIENT
// ============================================================================

class ApiClient {
  final Dio _dio;
  final FlutterSecureStorage _secureStorage;

  ApiClient({
    required String baseUrl,
    FlutterSecureStorage? secureStorage,
  })  : _dio = Dio(BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: const Duration(seconds: 30),
          receiveTimeout: const Duration(seconds: 30),
          sendTimeout: const Duration(seconds: 30),
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
          },
        )),
        _secureStorage = secureStorage ?? const FlutterSecureStorage() {
    _setupInterceptors();
  }

  Dio get dio => _dio;

  void _setupInterceptors() {
    // Auth Interceptor
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await _secureStorage.read(key: StorageKeys.accessToken);
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          return handler.next(options);
        },
        onError: (error, handler) async {
          if (error.response?.statusCode == 401) {
            // Token expired, try refresh
            final refreshed = await _refreshToken();
            if (refreshed) {
              // Retry the original request
              final opts = error.requestOptions;
              final token = await _secureStorage.read(key: StorageKeys.accessToken);
              opts.headers['Authorization'] = 'Bearer $token';
              
              try {
                final response = await _dio.fetch(opts);
                return handler.resolve(response);
              } catch (e) {
                return handler.next(error);
              }
            }
          }
          return handler.next(error);
        },
      ),
    );

    // Logging Interceptor (debug only)
    if (kDebugMode) {
      _dio.interceptors.add(
        LogInterceptor(
          requestHeader: true,
          requestBody: true,
          responseHeader: true,
          responseBody: true,
          error: true,
          logPrint: (log) => debugPrint(log.toString()),
        ),
      );
    }

    // Retry Interceptor
    _dio.interceptors.add(_RetryInterceptor(dio: _dio));
  }

  Future<bool> _refreshToken() async {
    try {
      final refreshToken = await _secureStorage.read(key: StorageKeys.refreshToken);
      if (refreshToken == null) return false;

      final response = await _dio.post(
        ApiEndpoints.refreshToken,
        data: {'refresh_token': refreshToken},
      );

      if (response.statusCode == 200) {
        final data = response.data;
        await _secureStorage.write(
          key: StorageKeys.accessToken,
          value: data['access_token'],
        );
        await _secureStorage.write(
          key: StorageKeys.refreshToken,
          value: data['refresh_token'],
        );
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  // ============================================================================
  // HTTP METHODS
  // ============================================================================

  Future<ApiResponse<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    try {
      final response = await _dio.get(
        path,
        queryParameters: queryParameters,
      );
      return ApiResponse.success(
        data: fromJson != null ? fromJson(response.data) : response.data,
        statusCode: response.statusCode ?? 200,
      );
    } on DioException catch (e) {
      return ApiResponse.error(
        message: _getErrorMessage(e),
        statusCode: e.response?.statusCode,
        error: e,
      );
    }
  }

  Future<ApiResponse<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    try {
      final response = await _dio.post(
        path,
        data: data,
        queryParameters: queryParameters,
      );
      return ApiResponse.success(
        data: fromJson != null ? fromJson(response.data) : response.data,
        statusCode: response.statusCode ?? 200,
      );
    } on DioException catch (e) {
      return ApiResponse.error(
        message: _getErrorMessage(e),
        statusCode: e.response?.statusCode,
        error: e,
      );
    }
  }

  Future<ApiResponse<T>> put<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    try {
      final response = await _dio.put(
        path,
        data: data,
        queryParameters: queryParameters,
      );
      return ApiResponse.success(
        data: fromJson != null ? fromJson(response.data) : response.data,
        statusCode: response.statusCode ?? 200,
      );
    } on DioException catch (e) {
      return ApiResponse.error(
        message: _getErrorMessage(e),
        statusCode: e.response?.statusCode,
        error: e,
      );
    }
  }

  Future<ApiResponse<T>> delete<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    T Function(Map<String, dynamic>)? fromJson,
  }) async {
    try {
      final response = await _dio.delete(
        path,
        queryParameters: queryParameters,
      );
      return ApiResponse.success(
        data: fromJson != null ? fromJson(response.data) : response.data,
        statusCode: response.statusCode ?? 200,
      );
    } on DioException catch (e) {
      return ApiResponse.error(
        message: _getErrorMessage(e),
        statusCode: e.response?.statusCode,
        error: e,
      );
    }
  }

  Future<ApiResponse<String>> uploadFile(
    String path,
    File file, {
    String fieldName = 'file',
    Map<String, dynamic>? extraData,
    void Function(int, int)? onProgress,
  }) async {
    try {
      final formData = FormData.fromMap({
        fieldName: await MultipartFile.fromFile(
          file.path,
          filename: file.path.split('/').last,
        ),
        ...?extraData,
      });

      final response = await _dio.post(
        path,
        data: formData,
        onSendProgress: onProgress,
      );

      return ApiResponse.success(
        data: response.data['url'] ?? response.data['id'],
        statusCode: response.statusCode ?? 200,
      );
    } on DioException catch (e) {
      return ApiResponse.error(
        message: _getErrorMessage(e),
        statusCode: e.response?.statusCode,
        error: e,
      );
    }
  }

  String _getErrorMessage(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return 'Connection timed out. Please check your internet connection.';
      case DioExceptionType.connectionError:
        return 'Unable to connect. Please check your internet connection.';
      case DioExceptionType.badResponse:
        final data = error.response?.data;
        if (data is Map && data.containsKey('message')) {
          return data['message'];
        }
        return 'Server error. Please try again later.';
      case DioExceptionType.cancel:
        return 'Request was cancelled.';
      default:
        return 'An unexpected error occurred.';
    }
  }
}

// ============================================================================
// RETRY INTERCEPTOR
// ============================================================================

class _RetryInterceptor extends Interceptor {
  final Dio dio;
  final int maxRetries;
  final List<int> retryStatusCodes;

  _RetryInterceptor({
    required this.dio,
    this.maxRetries = 3,
    this.retryStatusCodes = const [408, 500, 502, 503, 504],
  });

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    final extra = err.requestOptions.extra;
    final retryCount = (extra['retry_count'] as int?) ?? 0;

    if (retryCount < maxRetries &&
        (retryStatusCodes.contains(err.response?.statusCode) ||
            err.type == DioExceptionType.connectionTimeout ||
            err.type == DioExceptionType.sendTimeout ||
            err.type == DioExceptionType.receiveTimeout)) {
      await Future.delayed(Duration(seconds: retryCount + 1));
      
      final opts = err.requestOptions;
      opts.extra['retry_count'] = retryCount + 1;

      try {
        final response = await dio.fetch(opts);
        return handler.resolve(response);
      } catch (e) {
        // Continue to next retry or final error
      }
    }

    return handler.next(err);
  }
}

// ============================================================================
// API RESPONSE
// ============================================================================

class ApiResponse<T> {
  final T? data;
  final String? message;
  final int? statusCode;
  final bool isSuccess;
  final dynamic error;

  const ApiResponse._({
    this.data,
    this.message,
    this.statusCode,
    required this.isSuccess,
    this.error,
  });

  factory ApiResponse.success({
    T? data,
    String? message,
    int? statusCode,
  }) {
    return ApiResponse._(
      data: data,
      message: message,
      statusCode: statusCode,
      isSuccess: true,
    );
  }

  factory ApiResponse.error({
    String? message,
    int? statusCode,
    dynamic error,
  }) {
    return ApiResponse._(
      message: message,
      statusCode: statusCode,
      isSuccess: false,
      error: error,
    );
  }

  R when<R>({
    required R Function(T data) success,
    required R Function(String message) error,
  }) {
    if (isSuccess && data != null) {
      return success(data as T);
    } else {
      return error(message ?? 'An error occurred');
    }
  }
}

// ============================================================================
// PROVIDERS
// ============================================================================

final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient(baseUrl: ApiEndpoints.baseUrl);
});

final secureStorageProvider = Provider<FlutterSecureStorage>((ref) {
  return const FlutterSecureStorage();
});
