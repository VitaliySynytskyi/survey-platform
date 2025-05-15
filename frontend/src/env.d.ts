/// <reference types="vite/client" />

// Declare module for Vue SFC (Single File Components)
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// Extend Window interface to include custom properties
interface Window {
  // Add any custom window properties here if needed
}

// Declare global variables
declare const VITE_API_URL: string 