<script lang="ts">
	import { AddTask, UpdateTask } from '$lib/api/todo';
	import type { CreateTaskRequest, Task, Priority } from '$lib/types/todo';
	import { addTaskToStore, mapResponseToTask, updateInStore, tasks } from '$lib/stores/todo';
	import { notifications } from '$lib/stores/notifications';
	import { get } from 'svelte/store';

	let title = $state('');
	let description = $state('');
	let dueDate = $state('');
	let priority = $state<Priority>('MEDIUM');
	let category = $state('');
	let parentId = $state<string | undefined>(undefined);
	let editingTask: Task | null = $state(null);
	let minDate = $state(new Date().toISOString().split('T')[0]);

	let existingCategories = $derived(() => {
		const cats = new Set<string>();
		get(tasks).forEach((t) => {
			if (t.Category) cats.add(t.Category);
		});
		return Array.from(cats);
	});

	export function show(task: Task | null = null, forParentId?: string) {
		editingTask = task;
		parentId = forParentId;
		if (task) {
			title = task.Title;
			description = task.Description;
			dueDate = task.DueDate;
			priority = task.Priority || 'MEDIUM';
			category = task.Category || '';
		} else {
			title = '';
			description = '';
			dueDate = '';
			priority = 'MEDIUM';
			category = '';
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
				const payload = { title, description, dueDate, priority, category };
				const res = await UpdateTask(editingTask.Id, payload);
				const updated: Task = mapResponseToTask(res);
				updateInStore(updated);
				notifications.success('Task updated');
			} else {
				const taskReq: CreateTaskRequest = {
					title,
					description,
					dueDate,
					priority,
					category,
					parentId
				};
				const res = await AddTask(taskReq);
				const created = Array.isArray(res) ? res[0] : res;
				const taskForUi: Task = mapResponseToTask(created);
				addTaskToStore(taskForUi);
				notifications.success('Task created');
			}

			(form.closest('#add_modal') as HTMLDialogElement | null)?.close();
			form.reset();
			priority = 'MEDIUM';
			category = '';
			parentId = undefined;
		} catch (err) {
			notifications.error((err as Error).message);
		}
	}

	const priorityColors: Record<Priority, string> = {
		LOW: 'badge-ghost',
		MEDIUM: 'badge-info',
		HIGH: 'badge-warning',
		URGENT: 'badge-error'
	};
</script>

<dialog id="add_modal" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="mb-4 text-lg font-bold">
			{editingTask ? 'Edit Task' : parentId ? 'Add Subtask' : 'Create New Task'}
		</h3>
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
			<div class="flex gap-3">
				<label class="form-control w-1/2">
					<div class="label"><span class="label-text">Priority</span></div>
					<select class="select-bordered select w-full" bind:value={priority}>
						{#each ['LOW', 'MEDIUM', 'HIGH', 'URGENT'] as Priority[] as p}
							<option value={p}>{p}</option>
						{/each}
					</select>
				</label>
				<label class="form-control w-1/2">
					<div class="label"><span class="label-text">Category</span></div>
					<input
						type="text"
						list="category-list"
						class="input-bordered input w-full"
						placeholder="e.g. Work, Personal"
						bind:value={category}
					/>
					<datalist id="category-list">
						{#each existingCategories() as cat}
							<option value={cat}></option>
						{/each}
					</datalist>
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
				<button class="btn">Cancel</button>
			</form>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
