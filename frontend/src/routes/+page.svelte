<script lang="ts">
	import Fa from 'svelte-fa';
	import { faKey, faEnvelope, faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';

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

	async function postData(event: Event) {
		event.preventDefault();
		const url = 'http://localhost:9003/' + (isLoginPage ? 'login' : 'register');
		const reqData = { email: email, password: password };

		try {
			const res = await fetch(url, {
                method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify(reqData)
			});

			if (!res.ok) {
				throw new Error(`HTTP error! status: ${res.status}`);
			}

            if (!isLoginPage) {
                const data = await res.json();
                response = data.message;
                return
            }

            response = "user login successfully";
            console.log("response for login: ",res);
		} catch (error) {
			console.error('Error posting data:', error);
			response = `Error: ${error}`;
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
			<form class="flex flex-col items-center justify-center gap-4" onsubmit={postData}>
				<label class="input w-full">
					<Fa icon={faEnvelope}></Fa>
					<input
						bind:value={email}
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						class="grow"
						placeholder="e-mail"
					/>
				</label>
				<div class="w-full">
					<label class="input w-full">
						<Fa icon={faKey}></Fa>
						<input
							bind:value={password}
							id="password"
							name="password"
							type={getType()}
							required
							class="w-full grow"
							placeholder="password"
							minlength="8"
						/>
						<a
							href="#password"
							aria-label="password-eye"
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
				<button type="submit" class="btn btn-outline btn-primary lg:w-1/3">
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
					href="/"
					class="text-base-content hover:text-neutral font-semibold leading-6"
					onclick={() => {
						isLoginPage = !isLoginPage;
					}}
					>{#if isLoginPage}Register{:else}Login{/if}</a
				>
			</p>
		</div>
	</div>

	{#if response}
		<p>{response}</p>
	{/if}
</div>
