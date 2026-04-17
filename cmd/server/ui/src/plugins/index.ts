import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "@/App.vue";
import router from "@/router";
import { setupPermissionGuard } from "@/plugins/permission";
import "@/styles/global.css";
import "@/styles/console.scss";

export function bootstrap() {
  const app = createApp(App);
  app.use(createPinia());
  app.use(router);
  setupPermissionGuard(router);
  return app;
}
