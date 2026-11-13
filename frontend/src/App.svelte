<script lang="ts">
  import * as Tabs from "$lib/components/ui/tabs";
  import Create from "$lib/views/Create.svelte";
  import * as Separator from "$lib/components/ui/separator";
  import SSHKeys from "$lib/views/SSHKeys.svelte";
  import { ModeWatcher } from "mode-watcher";
</script>

<ModeWatcher defaultMode="dark" />

<div class="w-full min-w-0">
  <Tabs.Root value="#keys">
    <Tabs.List variant="line" class="ml-20">
      <Tabs.Trigger value="#keys">Keys</Tabs.Trigger>
      <Tabs.Trigger value="#create">Create</Tabs.Trigger>
    </Tabs.List>
    <Separator.Root class="my-2" />
    <Tabs.Content value="#keys">
      <SSHKeys />
    </Tabs.Content>
    <Tabs.Content value="#create">
      <Create />
    </Tabs.Content>
  </Tabs.Root>
</div>

<style>
  /* Put your standard CSS here */
</style>
