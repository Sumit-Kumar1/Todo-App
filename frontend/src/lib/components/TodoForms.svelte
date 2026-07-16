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
	let submitting = $state(false);

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

		if (!title || submitting) return;
		submitting = true;

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
		} finally {
			submitting = false;
		}
	}
</script>

<dialog id="add_modal" class="modal modal-bottom sm:modal-middle">
	<div class="modal-box w-full max-w-lg">
		<form method="dialog">
			<button class="btn absolute top-3 right-3 btn-circle btn-ghost btn-sm" aria-label="Close">
				<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		</form>

		<h3 class="mb-5 text-lg font-bold">
			{editingTask ? 'Edit Task' : parentId ? 'Add Subtask' : 'Create New Task'}
		</h3>

		<form class="flex flex-col gap-4" onsubmit={handleSubmit}>
			<label class="form-control w-full">
				<div class="label"><span class="label-text font-medium">Title</span></div>
				<input
					name="title"
					type="text"
					id="title"
					class="input-bordered input w-full"
					placeholder="What needs to be done?"
					required
					bind:value={title}
				/>
			</label>

			<label class="form-control w-full">
				<div class="label"><span class="label-text font-medium">Description</span></div>
				<input
					name="description"
					type="text"
					id="description"
					class="input-bordered input w-full"
					placeholder="Add some details (optional)"
					bind:value={description}
				/>
			</label>

			<label class="form-control w-full">
				<div class="label"><span class="label-text font-medium">Due Date</span></div>
				<input
					type="date"
					name="dueDate"
					id="dueDate"
					class="input-bordered input w-full"
					min={minDate}
					bind:value={dueDate}
				/>
			</label>

			<div class="flex flex-col gap-4 sm:flex-row sm:gap-3">
				<label class="form-control w-full sm:w-1/2">
					<div class="label"><span class="label-text font-medium">Priority</span></div>
					<select class="select-bordered select w-full" bind:value={priority}>
						{#each ['LOW', 'MEDIUM', 'HIGH', 'URGENT'] as Priority[] as p}
							<option value={p}>{p}</option>
						{/each}
					</select>
				</label>
				<label class="form-control w-full sm:w-1/2">
					<div class="label"><span class="label-text font-medium">Category</span></div>
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

			<div class="modal-action mt-2">
				<input type="reset" class="btn btn-ghost" value="Reset" />
				<button type="submit" class="btn btn-primary" disabled={submitting}>
					{#if submitting}
						<span class="loading loading-sm loading-spinner"></span>
					{/if}
					{editingTask ? 'Update Task' : 'Add Task'}
				</button>
			</div>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
