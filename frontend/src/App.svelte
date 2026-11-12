<script lang="ts">
  import NavigationBar from "$lib/views/components/NavigationBar.svelte";
  import Create from "$lib/views/Create.svelte";
  import SSHKeys from "$lib/views/SSHKeys.svelte";
  import { ModeWatcher } from "mode-watcher";

  let route = $state(window.location.hash || "#keys");

  function updateRoute() {
    route = window.location.hash || "#keys";
  }
</script>

<svelte:window onhashchange={updateRoute} />
<ModeWatcher defaultMode="dark" />

<div class="p-4 mt-12 w-full min-w-0">
  <NavigationBar />

  {#if route === "#keys"}
    <SSHKeys />
  {:else if route == "#create"}
    <Create />
  {:else}
    <p>Unknown Route: {route}</p>
  {/if}
</div>

<style>
  /* Put your standard CSS here */
</style>
