<script lang="ts">
	import { Fa } from 'svelte-fa';
	import { faEnvelope, faEye, faEyeSlash, faKey } from '@fortawesome/free-solid-svg-icons';
	import { Login, Register } from '$lib/api/user';
	import { goto } from '$app/navigation';
	import { notifications } from '$lib/stores/notifications';

	let passwordVisible = $state(false);
	let isLoginPage = $state(true);

	let email = $state('');
	let password = $state('');
	let response = $state('');

	let inputType = $derived(passwordVisible ? 'text' : 'password');

	async function submit(event: Event) {
		event.preventDefault();

		try {
			if (isLoginPage) {
				const res = await Login(email, password);
				if (res.status === 200) {
					notifications.success(res.data);
					goto('/todo');
				} else {
					notifications.error(res.error);
				}
			} else {
				const res = await Register(email, password);
				if (res.status === 201) {
					const loginRes = await Login(email, password);
					if (loginRes.data === 'user login successfully') {
						notifications.success(loginRes.data);
						goto('/todo');
					} else {
						notifications.error(loginRes.error);
					}
				} else {
					notifications.error(res.error);
				}
			}
		} catch (error) {
			const context = isLoginPage ? 'login' : 'register';
			console.error(`error while ${context}:`, error);
			notifications.error(`Unexpected ${context} error`);
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-base-200 text-base-content">
	<div
		class="card-border overflow-w-hidden card gap-2 border-base-300 bg-base-100 card-xl sm:w-2/3 lg:w-1/2"
	>
		<div class="card-title justify-center p-3">
			<h2 class="mt-5 text-center text-xl font-bold">
				{#if isLoginPage}
					Sign in to your account
				{:else}
					New User Registration
				{/if}
			</h2>
		</div>
		<div class="card-body gap-2">
			<form class="flex flex-col items-center justify-center gap-4" onsubmit={submit}>
				<label class="input w-full">
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
				<div class="w-full">
					<label class="input w-full">
						<Fa icon={faKey}></Fa>
						<input
							bind:value={password}
							class="w-full grow"
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
				</div>
				<button class="btn btn-outline btn-primary lg:w-1/3" type="submit">
					{#if isLoginPage}
						Sign in
					{:else}
						Register
					{/if}
				</button>
			</form>

			<p class="mt-5 text-center text-sm text-gray-500">
				{#if isLoginPage}
					Create new account?
				{:else}
					Already registered?
				{/if}
				<button
					type="button"
					class="btn leading-6 font-semibold text-base-content btn-link no-underline hover:text-neutral"
					onclick={() => {
						isLoginPage = !isLoginPage;
					}}
				>
					{#if isLoginPage}Register
					{:else}Login
					{/if}</button
				>
			</p>
		</div>
	</div>

	{#if response}
		<p>{response}</p>
	{/if}
</div>
