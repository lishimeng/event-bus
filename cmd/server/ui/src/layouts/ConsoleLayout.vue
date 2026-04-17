<script setup lang="ts">
import { SystemRouteName } from "@/enums/systemRoute";

const tabs: { to: string; name: string; label: string }[] = [
  { to: "/system/publish", name: SystemRouteName.Publish, label: "发布消息" },
  { to: "/system/records", name: SystemRouteName.Records, label: "消息记录" },
  { to: "/system/channel", name: SystemRouteName.Channel, label: "通道入库" },
  { to: "/system/rmq", name: SystemRouteName.Rmq, label: "RMQ 入库" },
];
</script>

<template>
  <div class="layout">
    <header class="hero">
      <div class="hero-inner">
        <p class="eyebrow">github.com/lishimeng/event-bus</p>
        <h1>Event Bus 控制台</h1>
        <p class="lede">
          构建产物在 <code>cmd/server/static</code>。<strong>通道入库</strong>写入
          <code>channel_config</code>；<strong>RMQ</strong> 写入
          <code>sys_config</code>（<code>name=rmq_config</code>）。改库后需<strong>重启服务</strong>才生效。
        </p>
      </div>
    </header>

    <nav class="tabs" aria-label="功能分区">
      <RouterLink v-for="t in tabs" :key="t.name" :to="t.to" class="tab">
        {{ t.label }}
      </RouterLink>
    </nav>

    <main class="panels">
      <keep-alive>
        <RouterView />
      </keep-alive>
    </main>

    <footer class="foot" />
  </div>
</template>
