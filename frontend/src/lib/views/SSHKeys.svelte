<script lang="ts">
  import SSHKeyTable from "./components/SSHKeyTable.svelte";
  import { SSHService, type SSHKey } from "bindings/quickssh/internal/ssh";
  import { onMount } from "svelte";

  let sshKeys: SSHKey[] = $state([]);

  onMount(() => {
    SSHService.GetExistingKeys()
      .then((keys) => {
        if (keys) {
          sshKeys = keys;
        }
      })
      .catch((e) => {
        console.log("Failed to get ssh keys: ", e);
      });
  });
</script>

<SSHKeyTable {sshKeys} />
