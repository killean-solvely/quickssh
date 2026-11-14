<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Field from "$lib/components/ui/field/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import * as Select from "$lib/components/ui/select";
  import {
    RemoteOS,
    SSHService,
    type DeployRequest,
    type HostKeyPrompt,
    type SSHKey,
  } from "bindings/quickssh/internal/ssh";

  interface Props {
    toDeployKey: SSHKey;
    closeDialog: () => void;
  }

  const { toDeployKey, closeDialog }: Props = $props();

  let host: string = $state("");
  let port = $state<number | undefined>(22);
  let username: string = $state("");
  let password: string = $state("");
  let remoteOS: RemoteOS = $state(RemoteOS.RemoteLinux);
  let busy = $state(false);
  let error = $state("");
  let approval = $state<{ request: DeployRequest; key: HostKeyPrompt } | null>(
    null,
  );

  const osTypes = [
    { value: RemoteOS.RemoteLinux, label: "Linux" },
    { value: RemoteOS.RemoteMacOS, label: "MacOS" },
    { value: RemoteOS.RemoteWindows, label: "Windows" },
  ];

  const canSubmit = $derived(
    host.trim().length > 0 &&
      port !== undefined &&
      Number.isInteger(port) &&
      port >= 1 &&
      port <= 65535 &&
      username.trim().length > 0 &&
      password.length > 0,
  );

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!canSubmit || busy || approval || port === undefined) return;

    await deploy({
      Host: host.trim(),
      Port: port,
      Username: username.trim(),
      Password: password,
      PublicKeyPath: toDeployKey.PublicKeyPath,
      RemoteMachineOS: remoteOS,
    });
  }

  async function deploy(request: DeployRequest) {
    busy = true;
    error = "";
    try {
      const result = await SSHService.DeployKey(request);
      if (result.hostKey) {
        approval = { request, key: result.hostKey };
        return;
      }
      host = "";
      port = 22;
      username = "";
      password = "";
      closeDialog();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      busy = false;
    }
  }

  async function trustAndRetry() {
    if (!approval || busy) return;
    const { request, key } = approval;
    approval = null;
    await deploy({ ...request, hostKeyApprovalID: key.approvalId });
  }
</script>

{#if approval}
  <section class="space-y-3" aria-labelledby="host-key-title">
    <h3 id="host-key-title" class="font-semibold">Trust this SSH host?</h3>
    <p class="text-sm">{approval.key.address} · {approval.key.algorithm}</p>
    <p class="select-all break-all text-sm">
      <code>{approval.key.fingerprint}</code>
    </p>
    <p class="text-muted-foreground text-sm">
      Compare this fingerprint with the one from the server console or your
      administrator.
    </p>
    <div class="flex justify-end gap-2">
      <Button type="button" variant="outline" onclick={() => (approval = null)}
        >Cancel</Button
      >
      <Button type="button" onclick={trustAndRetry}>Trust and retry</Button>
    </div>
  </section>
{:else}
  <form class="max-w-sm" onsubmit={handleSubmit}>
    <Field.Group>
      <Field.Field>
        <Field.Label for="host">Host</Field.Label>
        <Input id="host" placeholder="Host..." bind:value={host} />
      </Field.Field>
      <Field.Field>
        <Field.Label for="port">Port</Field.Label>
        <Input
          id="port"
          placeholder="Port..."
          type="number"
          min={1}
          max={65535}
          bind:value={port}
        />
      </Field.Field>
      <Field.Field>
        <Field.Label for="username">Username</Field.Label>
        <Input id="username" placeholder="Username..." bind:value={username} />
      </Field.Field>
      <Field.Field>
        <Field.Label for="password">Password</Field.Label>
        <Input
          id="password"
          type="password"
          placeholder="Password..."
          bind:value={password}
        />
      </Field.Field>
      <Field.Field>
        <Field.Label for="remoteos">Remote OS</Field.Label>
        <Select.Root type="single" items={osTypes} bind:value={remoteOS}>
          <Select.Trigger id="remoteos">
            <Select.Value />
          </Select.Trigger>
          <Select.Content>
            <Select.Group>
              {#each osTypes as osType (osType.value)}
                <Select.Item value={osType.value} label={osType.label}>
                  {osType.label}
                </Select.Item>
              {/each}
            </Select.Group>
          </Select.Content>
        </Select.Root>
      </Field.Field>
      <Field.Field orientation="horizontal">
        <Button type="submit" disabled={!canSubmit || busy}>
          {busy ? "Connecting…" : "Deploy Key"}
        </Button>
      </Field.Field>
    </Field.Group>
  </form>
{/if}

{#if error}
  <p class="text-destructive text-sm" role="alert">{error}</p>
{/if}
