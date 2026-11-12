<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Field from "$lib/components/ui/field/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import * as Select from "$lib/components/ui/select";
  import { SSHService } from "bindings/quickssh/internal/ssh";

  let keyName: string = $state("");
  let keyPassphrase: string = $state("");
  let keyComment: string = $state("");
  let keyType: string = $state("ed25519");

  const keyTypes = [
    { value: "ed25519", label: "ed25519" },
    { value: "rsa", label: "rsa" },
    { value: "ecdsa", label: "ecdsa" },
  ];

  const canSubmit = $derived(keyName.trim().length > 0);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!canSubmit) return;

    SSHService.CreateKey(keyName, keyPassphrase, keyType, keyComment)
      .then(() => {
        keyName = "";
        keyPassphrase = "";
        keyComment = "";
        keyType = "";
      })
      .catch((e) => {
        console.log("Something went wrong!", e);
      });
  }
</script>

<form class="max-w-sm" onsubmit={handleSubmit}>
  <Field.Group>
    <Field.Field>
      <Field.Label for="key-name">SSH Key Name</Field.Label>
      <Input id="key-name" placeholder="SSH Key Name..." bind:value={keyName} />
    </Field.Field>
    <Field.Field>
      <Field.Label for="key-passphrase">SSH Key Passphrase</Field.Label>
      <Input
        id="key-passphrase"
        placeholder="SSH Key Passphrase..."
        bind:value={keyPassphrase}
      />
      <Field.Description>Leave empty if not wanted</Field.Description>
    </Field.Field>
    <Field.Field>
      <Field.Label for="key-comment">SSH Key Comment</Field.Label>
      <Input
        id="key-comment"
        placeholder="SSH Key Comment..."
        bind:value={keyComment}
      />
    </Field.Field>
    <Field.Field>
      <Field.Label for="key-type">SSH Key Type</Field.Label>
      <Select.Root type="single" items={keyTypes} bind:value={keyType}>
        <Select.Trigger id="key-type">
          <Select.Value />
        </Select.Trigger>
        <Select.Content>
          <Select.Group>
            {#each keyTypes as kt (kt.value)}
              <Select.Item value={kt.value} label={kt.label}>
                {kt.label}
              </Select.Item>
            {/each}
          </Select.Group>
        </Select.Content>
      </Select.Root>
    </Field.Field>
    <Field.Field orientation="horizontal">
      <Button type="submit" disabled={!canSubmit}>Create Key</Button>
    </Field.Field>
  </Field.Group>
</form>
