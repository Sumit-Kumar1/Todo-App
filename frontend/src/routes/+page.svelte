<script lang="ts">
	import { Fa } from 'svelte-fa';
	import { faEnvelope, faEye, faEyeSlash, faKey } from '@fortawesome/free-solid-svg-icons';
	import { Login, Register } from '$lib/api/user';
	import { goto } from '$app/navigation';
	import { notifications } from '$lib/stores/notifications';

	let passwordVisible = $state(false);
	let isLoginPage = $state(true);
	let loading = $state(false);

	let email = $state('');
	let password = $state('');

	let inputType = $derived(passwordVisible ? 'text' : 'password');

	async function submit(event: Event) {
		event.preventDefault();
		if (loading) return;
		loading = true;

		try {
			if (isLoginPage) {
				const res = await Login(email, password);
				const status = res.status;

				if (status === 200) {
					notifications.success('user login success');
					goto('/todo');
				} else {
					notifications.error(res.text.toString());
				}
			} else {
				const res = await Register(email, password);
				const registerStatus = res.status;
				if (registerStatus === 201) {
					const loginRes = await Login(email, password);
					const loginText = loginRes.text.toString();

					if (loginText === 'user login successfully') {
						notifications.success('user login success');
						goto('/todo');
					} else {
						notifications.error(loginText);
					}
				} else {
					notifications.error(res.text.toString());
				}
			}
		} catch (error) {
			const errMsg = (error as Error).message;
			notifications.error(`${errMsg}`);
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-svh items-center justify-center bg-base-200 px-4 py-8 text-base-content">
	<div class="card w-full max-w-md border border-base-300 bg-base-100 shadow-lg">
		<div class="card-body gap-6">
			<div class="text-center">
				<h1 class="text-2xl font-bold tracking-tight">
					{#if isLoginPage}
						Welcome back
					{:else}
						Create an account
					{/if}
				</h1>
				<p class="mt-1 text-sm text-base-content/60">
					{#if isLoginPage}
						Sign in to manage your tasks
					{:else}
						Get started with your todo list
					{/if}
				</p>
			</div>

			<form class="flex flex-col gap-4" onsubmit={submit}>
				<label class="input-bordered input w-full">
					<Fa icon={faEnvelope}></Fa>
					<input
						autocomplete="email"
						bind:value={email}
						class="grow"
						id="email"
						name="email"
						placeholder="e-mail"
						required
						type="email"
					/>
				</label>

				<label class="input-bordered input w-full">
					<Fa icon={faKey}></Fa>
					<input
						bind:value={password}
						class="grow"
						id="password"
						minlength="8"
						name="password"
						placeholder="password"
						required
						type={inputType}
					/>
					<!-- svelte-ignore a11y_consider_explicit_label -->
					<button
						type="button"
						class="btn btn-ghost btn-xs"
						onclick={() => {
							passwordVisible = !passwordVisible;
						}}
					>
						{#if !passwordVisible}
							<Fa icon={faEye}></Fa>
						{:else}
							<Fa icon={faEyeSlash}></Fa>
						{/if}
					</button>
				</label>

				<button class="btn w-full btn-primary" type="submit" disabled={loading}>
					{#if loading}
						<span class="loading loading-sm loading-spinner"></span>
					{/if}
					{#if isLoginPage}
						Sign in
					{:else}
						Register
					{/if}
				</button>
			</form>

			<div class="divider my-0 text-xs text-base-content/40">OR</div>

			<p class="text-center text-sm text-base-content/60">
				{#if isLoginPage}
					Don't have an account?
				{:else}
					Already have an account?
				{/if}
				<button
					type="button"
					class="link font-semibold link-primary"
					onclick={() => {
						isLoginPage = !isLoginPage;
					}}
				>
					{#if isLoginPage}Register{:else}Login{/if}
				</button>
			</p>
		</div>
	</div>
</div>
