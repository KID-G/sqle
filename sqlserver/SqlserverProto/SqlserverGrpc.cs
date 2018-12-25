// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: Sqlserver.proto
// </auto-generated>
#pragma warning disable 0414, 1591
#region Designer generated code

using grpc = global::Grpc.Core;

namespace SqlserverProto {
  public static partial class SqlserverService
  {
    static readonly string __ServiceName = "SqlserverProto.SqlserverService";

    static readonly grpc::Marshaller<global::SqlserverProto.SplitSqlsInput> __Marshaller_SqlserverProto_SplitSqlsInput = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::SqlserverProto.SplitSqlsInput.Parser.ParseFrom);
    static readonly grpc::Marshaller<global::SqlserverProto.SplitSqlsOutput> __Marshaller_SqlserverProto_SplitSqlsOutput = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::SqlserverProto.SplitSqlsOutput.Parser.ParseFrom);
    static readonly grpc::Marshaller<global::SqlserverProto.AdviseInput> __Marshaller_SqlserverProto_AdviseInput = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::SqlserverProto.AdviseInput.Parser.ParseFrom);
    static readonly grpc::Marshaller<global::SqlserverProto.AdviseOutput> __Marshaller_SqlserverProto_AdviseOutput = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::SqlserverProto.AdviseOutput.Parser.ParseFrom);
    static readonly grpc::Marshaller<global::SqlserverProto.GetRollbackSqlsInput> __Marshaller_SqlserverProto_GetRollbackSqlsInput = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::SqlserverProto.GetRollbackSqlsInput.Parser.ParseFrom);
    static readonly grpc::Marshaller<global::SqlserverProto.GetRollbackSqlsOutput> __Marshaller_SqlserverProto_GetRollbackSqlsOutput = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::SqlserverProto.GetRollbackSqlsOutput.Parser.ParseFrom);

    static readonly grpc::Method<global::SqlserverProto.SplitSqlsInput, global::SqlserverProto.SplitSqlsOutput> __Method_GetSplitSqls = new grpc::Method<global::SqlserverProto.SplitSqlsInput, global::SqlserverProto.SplitSqlsOutput>(
        grpc::MethodType.Unary,
        __ServiceName,
        "GetSplitSqls",
        __Marshaller_SqlserverProto_SplitSqlsInput,
        __Marshaller_SqlserverProto_SplitSqlsOutput);

    static readonly grpc::Method<global::SqlserverProto.AdviseInput, global::SqlserverProto.AdviseOutput> __Method_Advise = new grpc::Method<global::SqlserverProto.AdviseInput, global::SqlserverProto.AdviseOutput>(
        grpc::MethodType.Unary,
        __ServiceName,
        "Advise",
        __Marshaller_SqlserverProto_AdviseInput,
        __Marshaller_SqlserverProto_AdviseOutput);

    static readonly grpc::Method<global::SqlserverProto.GetRollbackSqlsInput, global::SqlserverProto.GetRollbackSqlsOutput> __Method_GetRollbackSqls = new grpc::Method<global::SqlserverProto.GetRollbackSqlsInput, global::SqlserverProto.GetRollbackSqlsOutput>(
        grpc::MethodType.Unary,
        __ServiceName,
        "GetRollbackSqls",
        __Marshaller_SqlserverProto_GetRollbackSqlsInput,
        __Marshaller_SqlserverProto_GetRollbackSqlsOutput);

    /// <summary>Service descriptor</summary>
    public static global::Google.Protobuf.Reflection.ServiceDescriptor Descriptor
    {
      get { return global::SqlserverProto.SqlserverReflection.Descriptor.Services[0]; }
    }

    /// <summary>Base class for server-side implementations of SqlserverService</summary>
    public abstract partial class SqlserverServiceBase
    {
      public virtual global::System.Threading.Tasks.Task<global::SqlserverProto.SplitSqlsOutput> GetSplitSqls(global::SqlserverProto.SplitSqlsInput request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

      public virtual global::System.Threading.Tasks.Task<global::SqlserverProto.AdviseOutput> Advise(global::SqlserverProto.AdviseInput request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

      public virtual global::System.Threading.Tasks.Task<global::SqlserverProto.GetRollbackSqlsOutput> GetRollbackSqls(global::SqlserverProto.GetRollbackSqlsInput request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

    }

    /// <summary>Client for SqlserverService</summary>
    public partial class SqlserverServiceClient : grpc::ClientBase<SqlserverServiceClient>
    {
      /// <summary>Creates a new client for SqlserverService</summary>
      /// <param name="channel">The channel to use to make remote calls.</param>
      public SqlserverServiceClient(grpc::Channel channel) : base(channel)
      {
      }
      /// <summary>Creates a new client for SqlserverService that uses a custom <c>CallInvoker</c>.</summary>
      /// <param name="callInvoker">The callInvoker to use to make remote calls.</param>
      public SqlserverServiceClient(grpc::CallInvoker callInvoker) : base(callInvoker)
      {
      }
      /// <summary>Protected parameterless constructor to allow creation of test doubles.</summary>
      protected SqlserverServiceClient() : base()
      {
      }
      /// <summary>Protected constructor to allow creation of configured clients.</summary>
      /// <param name="configuration">The client configuration.</param>
      protected SqlserverServiceClient(ClientBaseConfiguration configuration) : base(configuration)
      {
      }

      public virtual global::SqlserverProto.SplitSqlsOutput GetSplitSqls(global::SqlserverProto.SplitSqlsInput request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return GetSplitSqls(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual global::SqlserverProto.SplitSqlsOutput GetSplitSqls(global::SqlserverProto.SplitSqlsInput request, grpc::CallOptions options)
      {
        return CallInvoker.BlockingUnaryCall(__Method_GetSplitSqls, null, options, request);
      }
      public virtual grpc::AsyncUnaryCall<global::SqlserverProto.SplitSqlsOutput> GetSplitSqlsAsync(global::SqlserverProto.SplitSqlsInput request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return GetSplitSqlsAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual grpc::AsyncUnaryCall<global::SqlserverProto.SplitSqlsOutput> GetSplitSqlsAsync(global::SqlserverProto.SplitSqlsInput request, grpc::CallOptions options)
      {
        return CallInvoker.AsyncUnaryCall(__Method_GetSplitSqls, null, options, request);
      }
      public virtual global::SqlserverProto.AdviseOutput Advise(global::SqlserverProto.AdviseInput request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return Advise(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual global::SqlserverProto.AdviseOutput Advise(global::SqlserverProto.AdviseInput request, grpc::CallOptions options)
      {
        return CallInvoker.BlockingUnaryCall(__Method_Advise, null, options, request);
      }
      public virtual grpc::AsyncUnaryCall<global::SqlserverProto.AdviseOutput> AdviseAsync(global::SqlserverProto.AdviseInput request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return AdviseAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual grpc::AsyncUnaryCall<global::SqlserverProto.AdviseOutput> AdviseAsync(global::SqlserverProto.AdviseInput request, grpc::CallOptions options)
      {
        return CallInvoker.AsyncUnaryCall(__Method_Advise, null, options, request);
      }
      public virtual global::SqlserverProto.GetRollbackSqlsOutput GetRollbackSqls(global::SqlserverProto.GetRollbackSqlsInput request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return GetRollbackSqls(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual global::SqlserverProto.GetRollbackSqlsOutput GetRollbackSqls(global::SqlserverProto.GetRollbackSqlsInput request, grpc::CallOptions options)
      {
        return CallInvoker.BlockingUnaryCall(__Method_GetRollbackSqls, null, options, request);
      }
      public virtual grpc::AsyncUnaryCall<global::SqlserverProto.GetRollbackSqlsOutput> GetRollbackSqlsAsync(global::SqlserverProto.GetRollbackSqlsInput request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return GetRollbackSqlsAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual grpc::AsyncUnaryCall<global::SqlserverProto.GetRollbackSqlsOutput> GetRollbackSqlsAsync(global::SqlserverProto.GetRollbackSqlsInput request, grpc::CallOptions options)
      {
        return CallInvoker.AsyncUnaryCall(__Method_GetRollbackSqls, null, options, request);
      }
      /// <summary>Creates a new instance of client from given <c>ClientBaseConfiguration</c>.</summary>
      protected override SqlserverServiceClient NewInstance(ClientBaseConfiguration configuration)
      {
        return new SqlserverServiceClient(configuration);
      }
    }

    /// <summary>Creates service definition that can be registered with a server</summary>
    /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
    public static grpc::ServerServiceDefinition BindService(SqlserverServiceBase serviceImpl)
    {
      return grpc::ServerServiceDefinition.CreateBuilder()
          .AddMethod(__Method_GetSplitSqls, serviceImpl.GetSplitSqls)
          .AddMethod(__Method_Advise, serviceImpl.Advise)
          .AddMethod(__Method_GetRollbackSqls, serviceImpl.GetRollbackSqls).Build();
    }

  }
}
#endregion