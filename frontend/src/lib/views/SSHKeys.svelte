<script lang="ts">
  import Badge from "$lib/components/ui/badge/badge.svelte";
  import * as Card from "$lib/components/ui/card";
  import KeyRound from "@lucide/svelte/icons/key-round";
  import X from "@lucide/svelte/icons/x";
  import Fingerprint from "@lucide/svelte/icons/fingerprint-pattern";
  import { SSHService, type SSHKey } from "bindings/quickssh/internal/ssh";
  import { onMount } from "svelte";
  import Separator from "$lib/components/ui/separator/separator.svelte";
  import * as Sheet from "$lib/components/ui/sheet";
  import { Button } from "$lib/components/ui/button";

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

  function deleteKey(sshKey: SSHKey) {
    SSHService.DeleteKey(sshKey).then(() => reloadKeys());
  }

  onMount(() => {
    reloadKeys();
  });
</script>

{#snippet sshKeyCard(k: SSHKey)}
  <Card.Root
    onclick={() => openKeyDetails(k)}
    class="w-full p-2 gap-0 hover:bg-accent/70"
  >
    <Card.Header class="flex align-middle px-1.5 pt-1.5 pb-2.5 justify-between">
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
    </Card.Header>
    <Card.Content class="px-2 pb-6">
      <div
        class="flex flex-row gap-3 align-middle text-muted-foreground text-xs"
      >
        <Fingerprint size="15px" class="self-center" />
        <p class="truncate">{k.Fingerprint.substring(0, 28)}...</p>
      </div>
    </Card.Content>
    <Card.Footer class="px-1.5 pb-1.5 flex flex-row align-middle">
      {#if k.InAgent}
        <div class="w-2.5 h-2.5 bg-green-500 rounded-full mr-2.5"></div>
        <p class="text-muted-foreground text-xs self-center">In Agent</p>
      {:else}
        <div class="w-2.5 h-2.5 bg-accent rounded-full mr-2.5"></div>
        <p class="text-muted-foreground text-xs self-center">Not in Agent</p>
      {/if}
    </Card.Footer>
  </Card.Root>
{/snippet}

<div class="px-4 flex h-full overflow-hidden">
  <div class="flex flex-col flex-1 min-w-0 overflow-y-auto">
    <div class="flex flex-col gap-1.5 mb-5 mt-4">
      <h1 class="text-3xl font-bold">SSH Keys</h1>
      <p class="text-muted-foreground">
        Manage your SSH identities and public keys.
      </p>
      <Separator class="mt-2.5" />
    </div>
    <div class="grid gap-4 grid-cols-1">
      {#each sshKeys as k}
        {@render sshKeyCard(k)}
      {/each}
    </div>
  </div>
  {#if selectedKey}
    <div class="flex shrink-0 overflow-hidden">
      <Separator orientation="vertical" class="ml-4" />
      <aside class="flex flex-col p-4">
        <div class="flex items-center justify-between">
          <h2 class="font-semibold">{selectedKey.Name}</h2>
          <Button variant="ghost" onclick={() => (selectedKey = null)}>
            <X />
          </Button>
        </div>

        <div class="min-h-0">
          <p class="break-all text-sm">
            {selectedKey.Fingerprint}
          </p>
        </div>
      </aside>
    </div>
  {/if}
</div>
