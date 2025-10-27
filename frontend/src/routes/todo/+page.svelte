<script lang="ts">
	import { Fa } from 'svelte-fa';
	import { faArrowRightFromBracket } from '@fortawesome/free-solid-svg-icons';
	import { Logout } from '$lib/api/user';
	import { goto } from '$app/navigation';
	import { notifications } from '$lib/stores/notifications';

	async function userLogout(event: Event) {
		event.preventDefault();

		try {
			const res = await Logout();
			if (res.data === 'user logged out successfully') {
				notifications.success('user logged out successfully');
			}
		} catch (err) {
			notifications.error((err as Error).message);
		}

		goto('/?page=login');
	}
</script>

<div>
	<div class="navbar border-b-2 border-accent p-2">
		<div class="flex-1">
			<p class="btn text-2xl btn-ghost">Todo App</p>
			<a href="/?page=api" class="btn btn-ghost">API Specification</a>
		</div>
		<div class="flex-none gap-2">
			<button class="btn btn-outline" aria-label="logout" onclick={userLogout}>
				<Fa icon={faArrowRightFromBracket}></Fa>
			</button>
		</div>
	</div>

	<div class="flex h-screen w-full flex-col items-center gap-5 p-3">
		<button
			class="btn w-1/3 btn-accent"
			aria-label="task-create"
			onclick={() => {
				// document.getElementById('add_modal').showModal();
			}}>Create New Task</button
		>
		<dialog id="add_modal" class="modal modal-bottom sm:modal-middle">
			<div class="modal-box">
				<form class="flex flex-col gap-3">
					<label class="floating-label">
						<input
							placeholder="Task name here..."
							name="title"
							type="text"
							id="title"
							class="validator input input-md w-full"
							required
							size="100"
						/>
						<span>Task name here...</span>
					</label>
					<label class="floating-label">
						<input
							placeholder="Description"
							name="description"
							type="text"
							id="description"
							class="validator input input-md w-full"
							size="1000"
						/>
						<span>Description</span>
					</label>
					<div>
						<label class="validator input">
							<span class="label">Due Date</span>
							<input
								type="date"
								name="dueDate"
								id="dueDate"
								required
								min="2025-01-01"
								max="2025-12-31"
							/>
						</label>
					</div>
					<div>
						<input type="reset" class="btn btn-outline btn-accent" />
						<button type="submit" class="btn btn-accent">Add Task</button>
					</div>
				</form>
				<div class="absolute right-3 bottom-4 modal-action p-2">
					<form method="dialog">
						<!-- if there is a button in form, it will close the modal -->
						<button class="btn">Cancel</button>
					</form>
				</div>
			</div>
		</dialog>
	</div>
</div>
