<script lang="ts">
  import Badge from "$lib/components/ui/badge/badge.svelte";
  import * as Card from "$lib/components/ui/card";
  import KeyRound from "@lucide/svelte/icons/key-round";
  import X from "@lucide/svelte/icons/x";
  import Fingerprint from "@lucide/svelte/icons/fingerprint-pattern";
  import { SSHService, type SSHKey } from "bindings/quickssh/internal/ssh";
  import { onMount } from "svelte";
  import Separator from "$lib/components/ui/separator/separator.svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Resizable from "$lib/components/ui/resizable";
  import { Clipboard } from "@wailsio/runtime";
  import Copy from "@lucide/svelte/icons/copy";
  import * as Table from "$lib/components/ui/table";

  let sshKeys: SSHKey[] = $state([]);
  let selectedKey = $state<SSHKey | null>(null);
  let sheetOpen = $state(false);

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

  function openKeyDetails(key: SSHKey) {
    selectedKey = key;
    sheetOpen = true;
  }

  function copyPublicKey(sshKey: SSHKey) {
    SSHService.GetPublicSSHKey(sshKey).then((val) => Clipboard.SetText(val));
  }

  function deleteKey(sshKey: SSHKey) {
    SSHService.DeleteKey(sshKey).then(() => reloadKeys());
  }

  onMount(() => {
    reloadKeys();
  });
</script>

{#snippet fingerprintSub(k: SSHKey)}
  <div class="flex flex-row gap-3 align-middle text-muted-foreground text-xs">
    <Fingerprint size="15px" class="self-center" />
    <p class="truncate">{k.Fingerprint.substring(0, 28)}...</p>
  </div>
{/snippet}

{#snippet sshTitleInfo(k: SSHKey)}
  <div class="flex align-middle">
    <div
      class="flex justify-center items-center w-12.5 h-12.5 bg-accent mr-4 rounded-lg text-muted-foreground"
    >
      <KeyRound />
    </div>
    <div class="self-center">
      <Card.Title>{k.Name}</Card.Title>
      <Card.Description>{k.Name}.pub</Card.Description>
    </div>
  </div>
  <Badge variant="outline" class="justify-self-end">{k.Algorithm}</Badge>
{/snippet}

{#snippet sshKeyActive(k: SSHKey)}
  {#if k.InAgent}
    <div class="w-2.5 h-2.5 bg-green-500 rounded-full mr-2.5"></div>
    <p class="text-muted-foreground text-xs self-center">In Agent</p>
  {:else}
    <div class="w-2.5 h-2.5 bg-accent rounded-full mr-2.5"></div>
    <p class="text-muted-foreground text-xs self-center">Not in Agent</p>
  {/if}
{/snippet}

{#snippet sshKeyCard(k: SSHKey)}
  <Card.Root
    onclick={() => openKeyDetails(k)}
    class={{
      "w-full p-2 gap-0 hover:bg-accent/70": true,
      "border-primary/50 border":
        selectedKey?.PublicKeyPath === k.PublicKeyPath,
    }}
  >
    <Card.Header class="flex align-middle px-1.5 pt-1.5 pb-2.5 justify-between">
      {@render sshTitleInfo(k)}
    </Card.Header>
    <Card.Content class="px-2 pb-6">
      {@render fingerprintSub(k)}
    </Card.Content>
    <Card.Footer class="px-1.5 pb-1.5 flex flex-row align-middle">
      {@render sshKeyActive(k)}
    </Card.Footer>
  </Card.Root>
{/snippet}

<Resizable.PaneGroup
  direction="horizontal"
  class="h-full min-h-0 overflow-hidden px-4"
>
  <Resizable.Pane class="min-w-0 min-h-0 overflow-hidden">
    <div class="flex h-full min-h-0 flex-col">
      <div class="mb-5 mt-4 flex shrink-0 flex-col gap-1.5">
        <h1 class="text-3xl font-bold">SSH Keys</h1>
        <p class="text-muted-foreground">
          Manage your SSH identities and public keys.
        </p>
        <Separator class="mt-2.5" />
      </div>

      <div class="min-h-0 flex-1 overflow-y-auto mb-3 rounded-xl">
        <div class="grid gap-4 grid-cols-1">
          {#each sshKeys as k}
            {@render sshKeyCard(k)}
          {/each}
        </div>
      </div>
    </div>
  </Resizable.Pane>
  {#if selectedKey}
    <Resizable.Handle withHandle class="ml-4" />
    <Resizable.Pane class="flex flex-col shrink-0 overflow-hidden px-4 pr-0">
      <aside class="flex w-full h-full min-h-0 flex-col">
        <div class="flex shrink-0 items-start justify-between py-4">
          {@render sshTitleInfo(selectedKey)}
        </div>
        <div class="min-h-0 flex-1">
          <Table.Root>
            <Table.Body>
              <Table.Row>
                <Table.Cell>Algorithm</Table.Cell>
                <Table.Cell>{selectedKey.Algorithm}</Table.Cell>
                <Table.Cell></Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>Fingerprint</Table.Cell>
                <Table.Cell>
                  {selectedKey.Fingerprint.substring(0, 28)}...
                </Table.Cell>
                <Table.Cell>
                  <Button
                    variant="ghost"
                    size="icon-xs"
                    class="p-0 m-0"
                    onclick={() => Clipboard.SetText(selectedKey!.Fingerprint)}
                  >
                    <Copy class="p-0 m-0" />
                  </Button>
                </Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>Comment</Table.Cell>
                <Table.Cell>{selectedKey.Comment}</Table.Cell>
                <Table.Cell></Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>Key Path</Table.Cell>
                <Table.Cell>{selectedKey.PrivateKeyPath}</Table.Cell>
                <Table.Cell>
                  <Button
                    variant="ghost"
                    size="icon-xs"
                    class="p-0 m-0"
                    onclick={() =>
                      Clipboard.SetText(
                        selectedKey!.PrivateKeyPath ??
                          selectedKey!.PublicKeyPath,
                      )}
                  >
                    <Copy class="p-0 m-0" />
                  </Button>
                </Table.Cell>
              </Table.Row>
              <Table.Row>
                <Table.Cell>Permissions</Table.Cell>
                <Table.Cell>{selectedKey.Permissions}</Table.Cell>
                <Table.Cell></Table.Cell>
              </Table.Row>
            </Table.Body>
          </Table.Root>
        </div>
        <div class="shrink-0 border-t pt-4">
          <Button
            class="w-full"
            onclick={() => copyPublicKey(selectedKey!)}
            disabled={!selectedKey}>Copy Public Key</Button
          >
          <Button
            variant="outline"
            class="w-full mb-2.5 mt-2"
            onclick={() => (selectedKey = null)}
          >
            Close
          </Button>
        </div>
      </aside>
    </Resizable.Pane>
  {/if}
</Resizable.PaneGroup>
