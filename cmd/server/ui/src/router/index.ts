import { createRouter, createWebHashHistory } from "vue-router";
import { SystemRouteName } from "@/enums/systemRoute";

const ConsoleLayout = () => import("@/layouts/ConsoleLayout.vue");

const router = createRouter({
  // 使用 hash 路由，避免服务端未做 history fallback 时刷新子路径 404。
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: ConsoleLayout,
      redirect: { name: SystemRouteName.Publish },
      children: [
        {
          path: "system/publish",
          name: SystemRouteName.Publish,
          meta: { title: "发布消息" },
          component: () => import("@/views/system/publish/PublishView.vue"),
        },
        {
          path: "system/records",
          name: SystemRouteName.Records,
          meta: { title: "消息记录" },
          component: () => import("@/views/system/records/RecordsView.vue"),
        },
        {
          path: "system/channel",
          name: SystemRouteName.Channel,
          meta: { title: "通道入库" },
          component: () => import("@/views/system/channel/ChannelView.vue"),
        },
        {
          path: "system/rmq",
          name: SystemRouteName.Rmq,
          meta: { title: "RMQ 入库" },
          component: () => import("@/views/system/rmq/RmqView.vue"),
        },
      ],
    },
  ],
});

router.afterEach((to) => {
  const t = to.meta.title as string | undefined;
  if (t) document.title = `${t} · Event Bus`;
});

export default router;
