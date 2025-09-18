import {
  type RequestHandlerOptions,
  type ResponseResolver,
  RequestHandler,
} from "msw";
import {
  type UniversalHandler,
  universalServerRequestFromFetch,
  universalServerResponseToFetch,
} from "@connectrpc/connect/protocol";
import {
  type ConnectRouter,
  type MethodImpl,
  type ServiceImpl,
  createConnectRouter,
} from "@connectrpc/connect";
import type { DescMethod, DescService } from "@bufbuild/protobuf";

export type ConnectHandlerInfo = {
  header: string;
  kind: "service" | "rpc";
  name: string;
};

type ConnectRequestParsedResult = {
  handler: UniversalHandler | undefined;
};

export type ConnectResolverExtras = {
  handler: UniversalHandler;
};

export type ConnectResponseResolver = ResponseResolver<ConnectResolverExtras>;

export type ConnectHandlerOptions = RequestHandlerOptions;
const defaultConnectResolver: ConnectResponseResolver = async ({
  request,
  handler,
}) => {
  const ureq = universalServerRequestFromFetch(request.clone(), {});
  const ures = await handler(ureq);
  return universalServerResponseToFetch(ures);
};

export class ConnectHandler extends RequestHandler<
  ConnectHandlerInfo,
  ConnectRequestParsedResult,
  ConnectResolverExtras
> {
  #router: ConnectRouter;
  constructor(
    info: Omit<ConnectHandlerInfo, "header">,
    routes: (router: ConnectRouter) => void,
    options?: ConnectHandlerOptions,
  ) {
    super({
      info: {
        ...info,
        header: `${info.kind} ${info.name}`,
      },
      resolver: defaultConnectResolver,
      options,
    });

    const router = createConnectRouter();
    routes(router);
    this.#router = router;
  }

  #handler(req: Request): UniversalHandler | undefined {
    const url = new URL(req.url);
    const handler = this.#router.handlers.find(
      (h) => h.requestPath == url.pathname,
    );
    if (handler == null) {
      return;
    }
    if (!handler.allowedMethods.includes(req.method)) {
      return;
    }
    return handler;
  }

  parse({ request }: { request: Request }) {
    return Promise.resolve({
      handler: this.#handler(request),
    });
  }

  predicate({
    parsedResult: { handler },
  }: {
    request: Request;
    parsedResult: ConnectRequestParsedResult;
  }) {
    if (handler == null) {
      return false;
    }
    return true;
  }

  extendResolverArgs({
    parsedResult: { handler },
  }: {
    request: Request;
    parsedResult: ConnectRequestParsedResult;
  }) {
    return {
      handler: handler!,
    };
  }

  log({
    request,
    response,
    parsedResult: { handler },
  }: {
    request: Request;
    response: Response;
    parsedResult: ConnectRequestParsedResult;
  }): void {
    if (handler == null) {
      throw new Error("handler is null");
    }
    const method = handler.method;

    console.groupCollapsed(
      `${method.methodKind} ${method.parent.typeName}/${method.name} (${response.status} ${response.statusText})`,
    );
    console.log("Request:", {
      url: new URL(request.url),
      method: request.method,
      headers: Object.fromEntries(request.headers.entries()),
    });
    console.log("Response:", {
      status: response.status,
      statusText: response.statusText,
      headers: Object.fromEntries(response.headers.entries()),
    });
    console.groupEnd();
  }
}

export const connect = {
  rpc: <M extends DescMethod>(
    method: M,
    impl: MethodImpl<M>,
    options?: ConnectHandlerOptions,
  ) =>
    new ConnectHandler(
      { kind: method.kind, name: `${method.parent.typeName}/${method.name}` },
      (router) => router.rpc(method, impl),
      options,
    ),
  service: <S extends DescService>(
    service: S,
    impl: ServiceImpl<S>,
    options?: ConnectHandlerOptions,
  ) =>
    new ConnectHandler(
      { kind: "service", name: service.typeName },
      (router) => router.service(service, impl),
      options,
    ),
};
