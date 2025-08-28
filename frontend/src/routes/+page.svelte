<script>
	import Fa from 'svelte-fa';
	import { faKey, faEnvelope, faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';

	let passwordVisible = $state.raw(false);
	let isLoginPage = $state.raw(false);

	function getType() {
		if (passwordVisible) {
			return 'text';
		} else {
			return 'password';
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
			<form class="flex flex-col items-center justify-center gap-4">
				<label class="input w-full">
					<Fa icon={faEnvelope}></Fa>
					<input
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
					class="leading-6 font-semibold text-base-content hover:text-neutral"
					onclick={() => {
						isLoginPage = !isLoginPage;
					}}
					>{#if isLoginPage}Register{:else}Login{/if}</a
				>
			</p>
		</div>
	</div>
</div>
