<script lang="ts">
  import { SSHService, type SSHKey } from "bindings/quickssh/internal/ssh";
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

<div class="p-4">
  <table>
    <tbody>
      {#each sshKeys as k}
        <tr>
          <td>{k.Path}</td>
          <td>{k.Name}</td>
          <td>{k.Pub}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  /* Put your standard CSS here */
</style>
