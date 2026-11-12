<script lang="ts">
  import * as ContextMenu from "$lib/components/ui/context-menu";
  import * as Table from "$lib/components/ui/table";
  import { Clipboard } from "@wailsio/runtime";
  import { SSHService, type SSHKey } from "bindings/quickssh/internal/ssh";
  import { onMount } from "svelte";

  let sshKeys: SSHKey[] = $state([]);

  function reloadKeys() {
    SSHService.GetExistingKeys()
      .then((keys) => {
        if (keys) {
          sshKeys = keys;
        }
      })
      .catch((e) => {
        console.log("Failed to get ssh keys: ", e);
      });
  }

  function deleteKey(sshKey: SSHKey) {
    SSHService.DeleteKey(sshKey.Path).then(() => reloadKeys());
  }

  onMount(() => {
    reloadKeys();
  });
</script>

<Table.Root class="w-full table-auto">
  <Table.Caption>Your public SSH keys</Table.Caption>
  <Table.Header>
    <Table.Row>
      <Table.Head class="w-px whitespace-nowrap">Name</Table.Head>
      <Table.Head>Public Key</Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#each sshKeys as k}
      <Table.Row>
        <Table.Cell class="w-px whitespace-nowrap">
          <ContextMenu.Root>
            <ContextMenu.Trigger>
              {k.Name}
            </ContextMenu.Trigger>
            <ContextMenu.Content>
              <ContextMenu.Item onclick={() => Clipboard.SetText(k.Path)}>
                Copy Path
              </ContextMenu.Item>
            </ContextMenu.Content>
          </ContextMenu.Root>
        </Table.Cell>
        <Table.Cell class="max-w-0">
          <ContextMenu.Root>
            <ContextMenu.Trigger class="block w-full truncate text-left">
              {k.Pub.length <= 20
                ? k.Pub
                : `${k.Pub.slice(0, 15)}...${k.Pub.slice(-15)}`}
            </ContextMenu.Trigger>
            <ContextMenu.Content>
              <ContextMenu.Item onclick={() => Clipboard.SetText(k.Pub)}>
                Copy
              </ContextMenu.Item>
              <ContextMenu.Separator />
              <ContextMenu.Item
                variant="destructive"
                onclick={() => deleteKey(k)}
              >
                Delete
              </ContextMenu.Item>
            </ContextMenu.Content>
          </ContextMenu.Root>
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
