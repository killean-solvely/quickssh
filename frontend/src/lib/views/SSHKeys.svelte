<script lang="ts">
  import Badge from "$lib/components/ui/badge/badge.svelte";
  import * as Card from "$lib/components/ui/card";
  import KeyRound from "@lucide/svelte/icons/key-round";
  import Fingerprint from "@lucide/svelte/icons/fingerprint-pattern";
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
    SSHService.DeleteKey(sshKey.PublicKeyPath).then(() => reloadKeys());
  }

  onMount(() => {
    reloadKeys();
  });
</script>

{#snippet sshKeyCard(k: SSHKey)}
  <Card.Root class="w-full p-2 gap-0 hover:bg-accent/70">
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
        <div class="w-[10px] h-[10px] bg-green-500 rounded-full mr-2.5"></div>
        <p class="text-muted-foreground text-xs self-center">In Agent</p>
      {:else}
        <div class="w-[10px] h-[10px] bg-accent rounded-full mr-2.5"></div>
        <p class="text-muted-foreground text-xs self-center">Not in Agent</p>
      {/if}
    </Card.Footer>
  </Card.Root>
{/snippet}

<div class="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 px-4">
  {#each sshKeys as k}
    {@render sshKeyCard(k)}
  {/each}
</div>
