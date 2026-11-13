<script lang="ts">
  import * as Tabs from "$lib/components/ui/tabs";
  import Create from "$lib/views/Create.svelte";
  import * as Separator from "$lib/components/ui/separator";
  import SSHKeys from "$lib/views/SSHKeys.svelte";
  import { ModeWatcher } from "mode-watcher";
  import KeyRound from "@lucide/svelte/icons/key-round";
</script>

<ModeWatcher defaultMode="dark" />

<div class="w-full min-w-0 h-dvh flex flex-col py-2">
  <Tabs.Root value="#keys" class="flex flex-col flex-1 min-h-0">
    <Tabs.List variant="line" class="ml-30 py-0 my-0">
      <Tabs.Trigger value="#keys">
        <div class="flex gap-2">
          <KeyRound class="self-center" />
          <p class="self-center">Keys</p>
        </div>
      </Tabs.Trigger>
      <Tabs.Trigger value="#create">Create</Tabs.Trigger>
    </Tabs.List>
    <Separator.Root />
    <Tabs.Content value="#keys" class="flex-1 min-h-0 overflow-hidden">
      <SSHKeys />
    </Tabs.Content>
    <Tabs.Content value="#create" class="flex-1 min-h-0 overflow-hidden">
      <Create />
    </Tabs.Content>
  </Tabs.Root>
</div>
