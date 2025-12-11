<script lang="ts">
	let minDate = $state(Date.now());
	import { AddTask } from '$lib/api/todo';
	import type { CreateTaskRequest, Task } from '$lib/types/todo';
	import { addTaskToStore, mapResponseToTask } from '$lib/stores/todo';
	import { notifications } from '$lib/stores/notifications';

	async function addTask(event: Event) {
		event.preventDefault();
		const form = event.target as HTMLFormElement;
		const formData = new FormData(form);
		const title = String(formData.get('title') ?? '').trim();
		const description = String(formData.get('description') ?? '').trim();
		const dueDate = String(formData.get('dueDate') ?? '').trim();

		if (!title) return;

		const taskReq: CreateTaskRequest = { title, description, dueDate };

		try {
			const res = await AddTask(taskReq);
			// If backend returns the created task, prefer it; otherwise use taskReq
			const created = Array.isArray(res) ? res[0] : res;
			const taskForUi: Task =
				created && (created as any).title
					? mapResponseToTask(created as any)
					: {
							Id: created.id,
							Title: created.title,
							Description: created.description,
							DueDate: dueDate,
							IsDone: false
						};
			addTaskToStore(taskForUi);
			(form.closest('#add_modal') as HTMLDialogElement | null)?.close();
			form.reset();
			notifications.success('Task created');
		} catch (err) {
			notifications.error((err as Error).message);
		}
	}
</script>

<dialog id="add_modal" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<form class="flex flex-col gap-3" onsubmit={addTask}>
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
				<span>Title...</span>
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
					<input type="date" name="dueDate" id="dueDate" min={minDate.toString()} />
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
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
