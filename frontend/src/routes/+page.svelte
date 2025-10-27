<script lang="ts">
	import { Fa } from 'svelte-fa';
	import { faEnvelope, faEye, faEyeSlash, faKey } from '@fortawesome/free-solid-svg-icons';
	import { Login, Register } from '$lib/api/user';
	import { goto } from '$app/navigation';

	let passwordVisible = $state.raw(false);
	let isLoginPage = $state.raw(true);

	let email = $state('');
	let password = $state('');
	let response = $state('');

	function getType() {
		if (passwordVisible) {
			return 'text';
		}

		return 'password';
	}

	async function submit(event: Event) {
		event.preventDefault();

		if (isLoginPage) {
			try {
				const res = await Login(email, password);
				goto('/todo');
			} catch (error) {
				console.error('error while login', error);
			}
			return;
		}

		try {
			await Register(email, password);
			goto('/');
		} catch (error) {
			console.error('error while register: ', error);
		}
	}
</script>

<div class="bg-base-200 text-base-content flex min-h-screen items-center justify-center">
	<div
		class="card-border overflow-w-hidden card border-base-300 bg-base-100 card-xl gap-2 sm:w-2/3 lg:w-1/2"
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
							type={getType()}
						/>
						<a
							aria-label="password-eye"
							href="#password"
							onclick={() => {
								passwordVisible = !passwordVisible;
							}}
						>
							{#if !passwordVisible}
								<Fa icon={faEye}></Fa>
							{:else}
								<Fa icon={faEyeSlash}></Fa>
							{/if}
						</a>
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
				<a
					class="text-base-content hover:text-neutral font-semibold leading-6"
					href="/"
					onclick={() => {
						isLoginPage = !isLoginPage;
					}}
				>
					{#if isLoginPage}Register
					{:else}Login
					{/if}</a
				>
			</p>
		</div>
	</div>

	{#if response}
		<p>{response}</p>
	{/if}
</div>
