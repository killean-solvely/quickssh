<script lang="ts">
  import * as NavigationMenu from "$lib/components/ui/navigation-menu";
  import * as Separator from "$lib/components/ui/separator";

  let route = $state(window.location.hash || "#keys");

  function updateRoute() {
    route = window.location.hash || "#keys";
  }
</script>

<svelte:window onhashchange={updateRoute} />

<NavigationMenu.Root>
  <NavigationMenu.List class="gap-2">
    <NavigationMenu.Item>
      <NavigationMenu.Link href="/#keys">Keys</NavigationMenu.Link>
    </NavigationMenu.Item>
    <NavigationMenu.Item>
      <NavigationMenu.Link href="/#create">Create</NavigationMenu.Link>
    </NavigationMenu.Item>
  </NavigationMenu.List>
</NavigationMenu.Root>

<Separator.Root class="my-2" />
