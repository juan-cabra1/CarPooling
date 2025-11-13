/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_USERS_API_URL: string;
  readonly VITE_TRIPS_API_URL: string;
  readonly VITE_BOOKINGS_API_URL: string;
  readonly VITE_SEARCH_API_URL: string;
  readonly VITE_APP_NAME: string;
  readonly VITE_APP_ENV: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
