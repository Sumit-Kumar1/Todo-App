<script lang="ts">
	import { AddTask, UpdateTask } from '$lib/api/todo';
	import type { CreateTaskRequest, Task } from '$lib/types/todo';
	import { addTaskToStore, mapResponseToTask, updateInStore } from '$lib/stores/todo';
	import { notifications } from '$lib/stores/notifications';

	let title = $state('');
	let description = $state('');
	let dueDate = $state('');
	let editingTask: Task | null = $state(null);
	let minDate = $state(new Date().toISOString().split('T')[0]);

	export function show(task: Task | null = null) {
		editingTask = task;
		if (task) {
			title = task.Title;
			description = task.Description;
			dueDate = task.DueDate;
		} else {
			title = '';
			description = '';
			dueDate = '';
		}
		const modal = document.getElementById('add_modal') as HTMLDialogElement | null;
		modal?.showModal();
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		const form = event.target as HTMLFormElement;

		if (!title) return;

		try {
			if (editingTask) {
				// Update existing task
				const payload = { title, description, dueDate };
				const res = await UpdateTask(editingTask.Id, payload);
				// Map response to store format if needed, or construct object
				// The API returns the updated task structure
				const updated: Task = {
					Id: res.id,
					Title: res.title,
					Description: res.description,
					DueDate: res.dueDate,
					IsDone: res.isDone
				};
				updateInStore(updated);
				notifications.success('Task updated');
			} else {
				// Create new task
				const taskReq: CreateTaskRequest = { title, description, dueDate };
				const res = await AddTask(taskReq);
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
				notifications.success('Task created');
			}

			(form.closest('#add_modal') as HTMLDialogElement | null)?.close();
			form.reset();
		} catch (err) {
			notifications.error((err as Error).message);
		}
	}
</script>

<dialog id="add_modal" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="mb-4 text-lg font-bold">{editingTask ? 'Edit Task' : 'Create New Task'}</h3>
		<form class="flex flex-col gap-3" onsubmit={handleSubmit}>
			<label class="floating-label">
				<input
					placeholder="Task name here..."
					name="title"
					type="text"
					id="title"
					class="validator input input-md w-full"
					required
					size="100"
					bind:value={title}
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
					bind:value={description}
				/>
				<span>Description</span>
			</label>
			<div>
				<label class="validator input">
					<span class="label">Due Date</span>
					<input type="date" name="dueDate" id="dueDate" min={minDate} bind:value={dueDate} />
				</label>
			</div>
			<div>
				<input type="reset" class="btn btn-outline btn-accent" value="Reset" />
				<button type="submit" class="btn btn-accent"
					>{editingTask ? 'Update Task' : 'Add Task'}</button
				>
			</div>
		</form>
		<div class="absolute right-3 bottom-0 modal-action p-2">
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
