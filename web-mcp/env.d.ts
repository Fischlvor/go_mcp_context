/// <reference types="vite/client" />

export interface ImportMetaEnv {
  VITE_SERVER_URL: string
  VITE_BASE_API: string
}

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
  }
}
