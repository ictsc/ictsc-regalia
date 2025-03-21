/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file was automatically generated by TanStack Router.
// You should NOT make any changes in this file as it will be overwritten.
// Additionally, you should also exclude this file from your linter and/or formatter to prevent it from being checked or modified.

// Import Routes

import { Route as rootRoute } from "./routes/~__root";
import { Route as TeamsRouteImport } from "./routes/~teams/~route";
import { Route as SignupRouteImport } from "./routes/~signup/~route";
import { Route as SigninRouteImport } from "./routes/~signin/~route";
import { Route as RuleRouteImport } from "./routes/~rule/~route";
import { Route as RankingRouteImport } from "./routes/~ranking/~route";
import { Route as ProblemsImport } from "./routes/~problems";
import { Route as AnnouncesImport } from "./routes/~announces";
import { Route as IndexRouteImport } from "./routes/~index/~route";
import { Route as ProblemsCodeRouteImport } from "./routes/~problems.$code/~route";
import { Route as AnnouncesSlugRouteImport } from "./routes/~announces.$slug/~route";
import { Route as ProblemsIndexRouteImport } from "./routes/~problems.index/~route";
import { Route as AnnouncesIndexRouteImport } from "./routes/~announces.index/~route";

// Create/Update Routes

const TeamsRouteRoute = TeamsRouteImport.update({
  id: "/teams",
  path: "/teams",
  getParentRoute: () => rootRoute,
} as any);

const SignupRouteRoute = SignupRouteImport.update({
  id: "/signup",
  path: "/signup",
  getParentRoute: () => rootRoute,
} as any);

const SigninRouteRoute = SigninRouteImport.update({
  id: "/signin",
  path: "/signin",
  getParentRoute: () => rootRoute,
} as any);

const RuleRouteRoute = RuleRouteImport.update({
  id: "/rule",
  path: "/rule",
  getParentRoute: () => rootRoute,
} as any).lazy(() => import("./routes/~rule/~route.lazy").then((d) => d.Route));

const RankingRouteRoute = RankingRouteImport.update({
  id: "/ranking",
  path: "/ranking",
  getParentRoute: () => rootRoute,
} as any);

const ProblemsRoute = ProblemsImport.update({
  id: "/problems",
  path: "/problems",
  getParentRoute: () => rootRoute,
} as any);

const AnnouncesRoute = AnnouncesImport.update({
  id: "/announces",
  path: "/announces",
  getParentRoute: () => rootRoute,
} as any);

const IndexRouteRoute = IndexRouteImport.update({
  id: "/",
  path: "/",
  getParentRoute: () => rootRoute,
} as any);

const ProblemsCodeRouteRoute = ProblemsCodeRouteImport.update({
  id: "/$code",
  path: "/$code",
  getParentRoute: () => ProblemsRoute,
} as any).lazy(() =>
  import("./routes/~problems.$code/~route.lazy").then((d) => d.Route),
);

const AnnouncesSlugRouteRoute = AnnouncesSlugRouteImport.update({
  id: "/$slug",
  path: "/$slug",
  getParentRoute: () => AnnouncesRoute,
} as any).lazy(() =>
  import("./routes/~announces.$slug/~route.lazy").then((d) => d.Route),
);

const ProblemsIndexRouteRoute = ProblemsIndexRouteImport.update({
  id: "/",
  path: "/",
  getParentRoute: () => ProblemsRoute,
} as any);

const AnnouncesIndexRouteRoute = AnnouncesIndexRouteImport.update({
  id: "/",
  path: "/",
  getParentRoute: () => AnnouncesRoute,
} as any);

// Populate the FileRoutesByPath interface

declare module "@tanstack/react-router" {
  interface FileRoutesByPath {
    "/": {
      id: "/";
      path: "/";
      fullPath: "/";
      preLoaderRoute: typeof IndexRouteImport;
      parentRoute: typeof rootRoute;
    };
    "/announces": {
      id: "/announces";
      path: "/announces";
      fullPath: "/announces";
      preLoaderRoute: typeof AnnouncesImport;
      parentRoute: typeof rootRoute;
    };
    "/problems": {
      id: "/problems";
      path: "/problems";
      fullPath: "/problems";
      preLoaderRoute: typeof ProblemsImport;
      parentRoute: typeof rootRoute;
    };
    "/ranking": {
      id: "/ranking";
      path: "/ranking";
      fullPath: "/ranking";
      preLoaderRoute: typeof RankingRouteImport;
      parentRoute: typeof rootRoute;
    };
    "/rule": {
      id: "/rule";
      path: "/rule";
      fullPath: "/rule";
      preLoaderRoute: typeof RuleRouteImport;
      parentRoute: typeof rootRoute;
    };
    "/signin": {
      id: "/signin";
      path: "/signin";
      fullPath: "/signin";
      preLoaderRoute: typeof SigninRouteImport;
      parentRoute: typeof rootRoute;
    };
    "/signup": {
      id: "/signup";
      path: "/signup";
      fullPath: "/signup";
      preLoaderRoute: typeof SignupRouteImport;
      parentRoute: typeof rootRoute;
    };
    "/teams": {
      id: "/teams";
      path: "/teams";
      fullPath: "/teams";
      preLoaderRoute: typeof TeamsRouteImport;
      parentRoute: typeof rootRoute;
    };
    "/announces/": {
      id: "/announces/";
      path: "/";
      fullPath: "/announces/";
      preLoaderRoute: typeof AnnouncesIndexRouteImport;
      parentRoute: typeof AnnouncesImport;
    };
    "/problems/": {
      id: "/problems/";
      path: "/";
      fullPath: "/problems/";
      preLoaderRoute: typeof ProblemsIndexRouteImport;
      parentRoute: typeof ProblemsImport;
    };
    "/announces/$slug": {
      id: "/announces/$slug";
      path: "/$slug";
      fullPath: "/announces/$slug";
      preLoaderRoute: typeof AnnouncesSlugRouteImport;
      parentRoute: typeof AnnouncesImport;
    };
    "/problems/$code": {
      id: "/problems/$code";
      path: "/$code";
      fullPath: "/problems/$code";
      preLoaderRoute: typeof ProblemsCodeRouteImport;
      parentRoute: typeof ProblemsImport;
    };
  }
}

// Create and export the route tree

interface AnnouncesRouteChildren {
  AnnouncesIndexRouteRoute: typeof AnnouncesIndexRouteRoute;
  AnnouncesSlugRouteRoute: typeof AnnouncesSlugRouteRoute;
}

const AnnouncesRouteChildren: AnnouncesRouteChildren = {
  AnnouncesIndexRouteRoute: AnnouncesIndexRouteRoute,
  AnnouncesSlugRouteRoute: AnnouncesSlugRouteRoute,
};

const AnnouncesRouteWithChildren = AnnouncesRoute._addFileChildren(
  AnnouncesRouteChildren,
);

interface ProblemsRouteChildren {
  ProblemsIndexRouteRoute: typeof ProblemsIndexRouteRoute;
  ProblemsCodeRouteRoute: typeof ProblemsCodeRouteRoute;
}

const ProblemsRouteChildren: ProblemsRouteChildren = {
  ProblemsIndexRouteRoute: ProblemsIndexRouteRoute,
  ProblemsCodeRouteRoute: ProblemsCodeRouteRoute,
};

const ProblemsRouteWithChildren = ProblemsRoute._addFileChildren(
  ProblemsRouteChildren,
);

export interface FileRoutesByFullPath {
  "/": typeof IndexRouteRoute;
  "/announces": typeof AnnouncesRouteWithChildren;
  "/problems": typeof ProblemsRouteWithChildren;
  "/ranking": typeof RankingRouteRoute;
  "/rule": typeof RuleRouteRoute;
  "/signin": typeof SigninRouteRoute;
  "/signup": typeof SignupRouteRoute;
  "/teams": typeof TeamsRouteRoute;
  "/announces/": typeof AnnouncesIndexRouteRoute;
  "/problems/": typeof ProblemsIndexRouteRoute;
  "/announces/$slug": typeof AnnouncesSlugRouteRoute;
  "/problems/$code": typeof ProblemsCodeRouteRoute;
}

export interface FileRoutesByTo {
  "/": typeof IndexRouteRoute;
  "/ranking": typeof RankingRouteRoute;
  "/rule": typeof RuleRouteRoute;
  "/signin": typeof SigninRouteRoute;
  "/signup": typeof SignupRouteRoute;
  "/teams": typeof TeamsRouteRoute;
  "/announces": typeof AnnouncesIndexRouteRoute;
  "/problems": typeof ProblemsIndexRouteRoute;
  "/announces/$slug": typeof AnnouncesSlugRouteRoute;
  "/problems/$code": typeof ProblemsCodeRouteRoute;
}

export interface FileRoutesById {
  __root__: typeof rootRoute;
  "/": typeof IndexRouteRoute;
  "/announces": typeof AnnouncesRouteWithChildren;
  "/problems": typeof ProblemsRouteWithChildren;
  "/ranking": typeof RankingRouteRoute;
  "/rule": typeof RuleRouteRoute;
  "/signin": typeof SigninRouteRoute;
  "/signup": typeof SignupRouteRoute;
  "/teams": typeof TeamsRouteRoute;
  "/announces/": typeof AnnouncesIndexRouteRoute;
  "/problems/": typeof ProblemsIndexRouteRoute;
  "/announces/$slug": typeof AnnouncesSlugRouteRoute;
  "/problems/$code": typeof ProblemsCodeRouteRoute;
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath;
  fullPaths:
    | "/"
    | "/announces"
    | "/problems"
    | "/ranking"
    | "/rule"
    | "/signin"
    | "/signup"
    | "/teams"
    | "/announces/"
    | "/problems/"
    | "/announces/$slug"
    | "/problems/$code";
  fileRoutesByTo: FileRoutesByTo;
  to:
    | "/"
    | "/ranking"
    | "/rule"
    | "/signin"
    | "/signup"
    | "/teams"
    | "/announces"
    | "/problems"
    | "/announces/$slug"
    | "/problems/$code";
  id:
    | "__root__"
    | "/"
    | "/announces"
    | "/problems"
    | "/ranking"
    | "/rule"
    | "/signin"
    | "/signup"
    | "/teams"
    | "/announces/"
    | "/problems/"
    | "/announces/$slug"
    | "/problems/$code";
  fileRoutesById: FileRoutesById;
}

export interface RootRouteChildren {
  IndexRouteRoute: typeof IndexRouteRoute;
  AnnouncesRoute: typeof AnnouncesRouteWithChildren;
  ProblemsRoute: typeof ProblemsRouteWithChildren;
  RankingRouteRoute: typeof RankingRouteRoute;
  RuleRouteRoute: typeof RuleRouteRoute;
  SigninRouteRoute: typeof SigninRouteRoute;
  SignupRouteRoute: typeof SignupRouteRoute;
  TeamsRouteRoute: typeof TeamsRouteRoute;
}

const rootRouteChildren: RootRouteChildren = {
  IndexRouteRoute: IndexRouteRoute,
  AnnouncesRoute: AnnouncesRouteWithChildren,
  ProblemsRoute: ProblemsRouteWithChildren,
  RankingRouteRoute: RankingRouteRoute,
  RuleRouteRoute: RuleRouteRoute,
  SigninRouteRoute: SigninRouteRoute,
  SignupRouteRoute: SignupRouteRoute,
  TeamsRouteRoute: TeamsRouteRoute,
};

export const routeTree = rootRoute
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>();

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "~__root.tsx",
      "children": [
        "/",
        "/announces",
        "/problems",
        "/ranking",
        "/rule",
        "/signin",
        "/signup",
        "/teams"
      ]
    },
    "/": {
      "filePath": "~index/~route.tsx"
    },
    "/announces": {
      "filePath": "~announces.tsx",
      "children": [
        "/announces/",
        "/announces/$slug"
      ]
    },
    "/problems": {
      "filePath": "~problems.tsx",
      "children": [
        "/problems/",
        "/problems/$code"
      ]
    },
    "/ranking": {
      "filePath": "~ranking/~route.tsx"
    },
    "/rule": {
      "filePath": "~rule/~route.tsx"
    },
    "/signin": {
      "filePath": "~signin/~route.tsx"
    },
    "/signup": {
      "filePath": "~signup/~route.tsx"
    },
    "/teams": {
      "filePath": "~teams/~route.tsx"
    },
    "/announces/": {
      "filePath": "~announces.index/~route.tsx",
      "parent": "/announces"
    },
    "/problems/": {
      "filePath": "~problems.index/~route.tsx",
      "parent": "/problems"
    },
    "/announces/$slug": {
      "filePath": "~announces.$slug/~route.tsx",
      "parent": "/announces"
    },
    "/problems/$code": {
      "filePath": "~problems.$code/~route.tsx",
      "parent": "/problems"
    }
  }
}
ROUTE_MANIFEST_END */
