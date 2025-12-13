<script lang="ts">
	export const prerender = false;
	import { Fa } from 'svelte-fa';
	import { faArrowRightFromBracket } from '@fortawesome/free-solid-svg-icons';
	import { Logout } from '$lib/api/user';
	import { goto } from '$app/navigation';
	import { notifications } from '$lib/stores/notifications';
	import TodoForms from '$lib/components/TodoForms.svelte';
	import { onMount } from 'svelte';
	import {
		deletedTasks,
		loadTasks,
		markDoneInStore,
		moveToDeleted,
		tasks,
	} from '$lib/stores/todo';
	import { DelTask, MarkDone } from '$lib/api/todo';
	import type { Task } from '$lib/types/todo';

	async function userLogout(event: Event) {
		event.preventDefault();

		try {
			const res = await Logout();
			const status = res.status;
			if (status == 200 && res.text.toString() === 'user logged out successfully') {
				notifications.success('user logged out successfully');
			}
		} catch (err) {
			notifications.error((err as Error).message, 5000);
		}

		goto('/');
	}

	let todoForm: TodoForms;

	function openCreateModal() {
		todoForm.show(null);
	}

	onMount(() => {
		loadTasks();
	});

	let activeTab: 'tasks' | 'done' = 'tasks';

	async function handleDelete(id: string) {
		try {
			await DelTask(id);
			moveToDeleted(id);
		} catch (err) {
			console.error((err as Error).message);
		}
	}

	async function handleMarkDone(id: string) {
		try {
			const res = await MarkDone(id);
			if (res) {
				markDoneInStore(id);
			}
		} catch (err) {
			console.error((err as Error).message);
		}
	}

	async function handleEdit(task: Task) {
		todoForm.show(task);
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
		<button class="btn w-1/3 btn-accent" aria-label="task-create" onclick={openCreateModal}
			>Create New Task</button
		>

		<TodoForms bind:this={todoForm} />

		<div class="w-full max-w-3xl">
			<div role="tablist" class="tabs-lift mb-3 tabs justify-center">
				<button
					role="tab"
					class={`tab ${activeTab === 'tasks' ? 'tab-active' : ''}`}
					onclick={() => (activeTab = 'tasks')}>Tasks</button
				>
				<button
					role="tab"
					class={`tab ${activeTab === 'done' ? 'tab-active' : ''}`}
					onclick={() => (activeTab = 'done')}>Done</button
				>
			</div>
			{#if $tasks.length === 0}
				<p class="text-center opacity-70">No tasks yet.</p>
			{:else}
				<ul class="flex flex-col gap-2">
					{#if activeTab === 'tasks'}
						{#each $tasks.filter((t) => !t.IsDone) as t}
							<li class="card bg-base-200 p-4">
								<div class="flex items-start justify-between gap-3">
									<div>
										<p class="font-semibold">{t.Title}</p>
										<p class="text-sm opacity-80">{t.Description}</p>
										<p class="text-xs opacity-60">Due: {t.DueDate}</p>
									</div>
									<div class="flex gap-2">
										<button
											class="btn btn-ghost btn-sm"
											title="Edit"
											aria-label="edit button"
											onclick={() => handleEdit(t)}
										>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
												/>
											</svg>
										</button>
										<button
											class="btn btn-ghost btn-sm"
											title="Delete"
											aria-label="Delete button"
											onclick={() => handleDelete(t.Id)}
										>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
												/>
											</svg>
										</button>
										<button
											class="btn btn-ghost btn-sm"
											title="Mark done"
											aria-label="Mark done button"
											onclick={() => handleMarkDone(t.Id)}
										>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M5 13l4 4L19 7"
												/>
											</svg>
										</button>
									</div>
								</div>
							</li>
						{/each}
					{:else if activeTab === 'done'}
						{#each $tasks.filter((t) => t.IsDone) as t}
							<li class="card bg-base-200 p-4">
								<div class="flex items-start justify-between gap-3">
									<div>
										<p class="font-semibold">{t.Title}</p>
										<p class="text-sm opacity-80">{t.Description}</p>
										<p class="text-xs opacity-60">Due: {t.DueDate}</p>
										<p class="text-xs text-success">Done</p>
									</div>
									<div class="flex gap-2">
										<button
											class="btn btn-ghost btn-sm"
											aria-label="delete"
											title="Delete"
											onclick={() => handleDelete(t.Id)}
										>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
												/>
											</svg>
										</button>
									</div>
								</div>
							</li>
						{/each}
					{:else}
						{#each $deletedTasks as t}
							<li class="card bg-base-200 p-4">
								<div class="flex items-start justify-between gap-3">
									<div>
										<p class="font-semibold">{t.Title}</p>
										<p class="text-sm opacity-80">{t.Description}</p>
										<p class="text-xs opacity-60">Due: {t.DueDate}</p>
									</div>
									<div class="flex gap-2 text-xl">
										<!-- No edit or mark-done for deleted -->
									</div>
								</div>
							</li>
						{/each}
					{/if}
				</ul>
			{/if}
		</div>
	</div>
</div>
