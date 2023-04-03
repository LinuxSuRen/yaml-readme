/**
 * @fileoverview gRPC-Web generated client stub for server
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.4.2
// 	protoc              v3.12.4
// source: server.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.server = require('./server_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.RunnerClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.server.RunnerPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.server.TestTask,
 *   !proto.server.HelloReply>}
 */
const methodDescriptor_Runner_Run = new grpc.web.MethodDescriptor(
  '/server.Runner/Run',
  grpc.web.MethodType.UNARY,
  proto.server.TestTask,
  proto.server.HelloReply,
  /**
   * @param {!proto.server.TestTask} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.server.HelloReply.deserializeBinary
);


/**
 * @param {!proto.server.TestTask} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.server.HelloReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.server.HelloReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.server.RunnerClient.prototype.run =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/server.Runner/Run',
      request,
      metadata || {},
      methodDescriptor_Runner_Run,
      callback);
};


/**
 * @param {!proto.server.TestTask} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.server.HelloReply>}
 *     Promise that resolves to the response
 */
proto.server.RunnerPromiseClient.prototype.run =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/server.Runner/Run',
      request,
      metadata || {},
      methodDescriptor_Runner_Run);
};


module.exports = proto.server;

