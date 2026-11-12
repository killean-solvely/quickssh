<script lang="ts">
  import * as ContextMenu from "$lib/components/ui/context-menu";
  import * as Table from "$lib/components/ui/table";
  import { Clipboard } from "@wailsio/runtime";
  import type { SSHKey } from "bindings/quickssh/internal/ssh";

  interface Props {
    sshKeys: SSHKey[];
  }

  const { sshKeys }: Props = $props();
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
              {k.Pub}
            </ContextMenu.Trigger>
            <ContextMenu.Content>
              <ContextMenu.Item onclick={() => Clipboard.SetText(k.Pub)}>
                Copy
              </ContextMenu.Item>
            </ContextMenu.Content>
          </ContextMenu.Root>
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
