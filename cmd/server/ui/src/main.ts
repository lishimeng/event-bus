import { bootstrap } from "@/plugins";
import { setupDirectives } from "@/directive";

const app = bootstrap();
setupDirectives(app);
app.mount("#app");
