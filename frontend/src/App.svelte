<script lang="ts">
  import * as Table from "$lib/components/ui/table";
  import { SSHService, type SSHKey } from "bindings/quickssh/internal/ssh";
  import { ModeWatcher } from "mode-watcher";
  import { onMount } from "svelte";

  let sshKeys: SSHKey[] = $state([]);

  onMount(() => {
    SSHService.GetExistingKeys().then((keys) => {
      if (keys) {
        sshKeys = keys;
      }
    });
  });
</script>

<ModeWatcher defaultMode="dark" />

<div class="p-4">
  <Table.Root>
    <Table.Caption>A list of public SSH keys</Table.Caption>
    <Table.Header>
      <Table.Row>
        <Table.Head>Path</Table.Head>
        <Table.Head>Name</Table.Head>
        <Table.Head>Public Key</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#each sshKeys as k}
        <Table.Row>
          <Table.Cell>{k.Path}</Table.Cell>
          <Table.Cell>{k.Name}</Table.Cell>
          <Table.Cell>{k.Pub}</Table.Cell>
        </Table.Row>
      {/each}
    </Table.Body>
  </Table.Root>
</div>

<style>
  /* Put your standard CSS here */
</style>
