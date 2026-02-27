<script lang="ts">
	export const prerender = false;
	import { Fa } from 'svelte-fa';
	import { faArrowRightFromBracket, faSun, faMoon } from '@fortawesome/free-solid-svg-icons';
	import { Logout } from '$lib/api/user';
	import { goto } from '$app/navigation';
	import { notifications } from '$lib/stores/notifications';
	import TodoForms from '$lib/components/TodoForms.svelte';
	import { onMount } from 'svelte';
	import {
		clearCompleted,
		loadTasks,
		markDoneInStore,
		moveToDeleted,
		tasks
	} from '$lib/stores/todo';
	import { DelTask, MarkDone } from '$lib/api/todo';
	import type { Task, Priority } from '$lib/types/todo';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { dndzone } from 'svelte-dnd-action';
	import { browser } from '$app/environment';

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
		loadDragOrder();
	});

	let activeTab: 'all' | 'active' | 'done' = $state('all');
	let searchQuery = $state('');
	let filterPriority = $state('');
	let filterCategory = $state('');
	let searchTimeout: ReturnType<typeof setTimeout>;
	let debouncedSearch = $state('');
	let expandedTasks = $state<Set<string>>(new Set());

	function onSearchInput(e: Event) {
		const val = (e.target as HTMLInputElement).value;
		searchQuery = val;
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			debouncedSearch = val;
		}, 300);
	}

	let allCategories = $derived(() => {
		const cats = new Set<string>();
		$tasks.forEach((t) => {
			if (t.Category) cats.add(t.Category);
		});
		return Array.from(cats).sort();
	});

	// Build parent tasks with children nested
	let parentTasks = $derived(() => {
		const parents = $tasks.filter((t) => !t.ParentId);
		return parents.map((p) => ({
			...p,
			ChildTasks: $tasks.filter((c) => c.ParentId === p.Id)
		}));
	});

	let filteredTasks = $derived(() => {
		let list = parentTasks();
		const search = debouncedSearch.toLowerCase();

		if (activeTab === 'active') {
			list = list.filter((t) => !t.IsDone);
		} else if (activeTab === 'done') {
			list = list.filter((t) => t.IsDone);
		}

		if (filterPriority) {
			list = list.filter((t) => t.Priority === filterPriority);
		}

		if (filterCategory) {
			list = list.filter((t) => t.Category === filterCategory);
		}

		if (search) {
			list = list.filter(
				(t) =>
					t.Title.toLowerCase().includes(search) || t.Description.toLowerCase().includes(search)
			);
		}

		// Apply saved drag order
		if (dragOrder.length > 0) {
			const orderMap = new Map(dragOrder.map((id, idx) => [id, idx]));
			list = [...list].sort((a, b) => {
				const ai = orderMap.get(a.Id) ?? Infinity;
				const bi = orderMap.get(b.Id) ?? Infinity;
				return ai - bi;
			});
		}

		return list;
	});

	let doneCount = $derived($tasks.filter((t) => !t.ParentId && t.IsDone).length);

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

	function handleAddSubtask(parentId: string) {
		todoForm.show(null, parentId);
	}

	function toggleExpand(id: string) {
		expandedTasks = new Set(expandedTasks);
		if (expandedTasks.has(id)) {
			expandedTasks.delete(id);
		} else {
			expandedTasks.add(id);
		}
	}

	async function handleClearCompleted() {
		try {
			await clearCompleted();
			notifications.success('Completed tasks cleared');
		} catch (err) {
			notifications.error((err as Error).message);
		}
	}

	// Drag and drop
	let dragOrder: string[] = $state([]);

	function loadDragOrder() {
		if (browser) {
			try {
				const saved = localStorage.getItem('taskOrder');
				if (saved) dragOrder = JSON.parse(saved);
			} catch {
				/* ignore */
			}
		}
	}

	function saveDragOrder(items: Task[]) {
		dragOrder = items.map((t) => t.Id);
		if (browser) {
			localStorage.setItem('taskOrder', JSON.stringify(dragOrder));
		}
	}

	type DndTask = Task & { id: string };

	function handleDndConsider(e: CustomEvent<{ items: DndTask[] }>) {
		dndItems = e.detail.items;
	}

	function handleDndFinalize(e: CustomEvent<{ items: DndTask[] }>) {
		dndItems = e.detail.items;
		saveDragOrder(dndItems);
	}

	// Wrap filtered tasks for dnd (needs id property)
	let dndItems: DndTask[] = $state([]);

	$effect(() => {
		dndItems = filteredTasks().map((t) => ({
			...t,
			id: t.Id
		}));
	});

	const priorityBorder: Record<Priority, string> = {
		LOW: 'border-l-gray-400',
		MEDIUM: 'border-l-info',
		HIGH: 'border-l-warning',
		URGENT: 'border-l-error'
	};

	const priorityBadge: Record<Priority, string> = {
		LOW: 'badge-ghost',
		MEDIUM: 'badge-info',
		HIGH: 'badge-warning',
		URGENT: 'badge-error'
	};
</script>

<div>
	<div class="navbar border-b-2 border-accent p-2">
		<div class="flex-1">
			<p class="btn text-2xl btn-ghost">Todo App</p>
			<a href="/?page=api" class="btn btn-ghost">API Specification</a>
		</div>
		<div class="flex-none gap-2">
			<button class="btn btn-circle btn-ghost" aria-label="toggle theme" onclick={toggleTheme}>
				{#if $theme === 'light'}
					<Fa icon={faMoon} />
				{:else}
					<Fa icon={faSun} />
				{/if}
			</button>
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
			<!-- Search bar -->
			<div class="mb-3">
				<input
					type="text"
					class="input-bordered input w-full"
					placeholder="Search tasks..."
					value={searchQuery}
					oninput={onSearchInput}
				/>
			</div>

			<!-- Filter row -->
			<div class="mb-3 flex flex-wrap items-center gap-3">
				<!-- Tab filters -->
				<div role="tablist" class="tabs-lift tabs">
					<button
						role="tab"
						class={`tab ${activeTab === 'all' ? 'tab-active' : ''}`}
						onclick={() => (activeTab = 'all')}>All</button
					>
					<button
						role="tab"
						class={`tab ${activeTab === 'active' ? 'tab-active' : ''}`}
						onclick={() => (activeTab = 'active')}>Active</button
					>
					<button
						role="tab"
						class={`tab ${activeTab === 'done' ? 'tab-active' : ''}`}
						onclick={() => (activeTab = 'done')}>Done</button
					>
				</div>

				<!-- Priority filter -->
				<select class="select-bordered select select-sm" bind:value={filterPriority}>
					<option value="">All Priorities</option>
					<option value="URGENT">Urgent</option>
					<option value="HIGH">High</option>
					<option value="MEDIUM">Medium</option>
					<option value="LOW">Low</option>
				</select>

				<!-- Category filter -->
				<select class="select-bordered select select-sm" bind:value={filterCategory}>
					<option value="">All Categories</option>
					{#each allCategories() as cat}
						<option value={cat}>{cat}</option>
					{/each}
				</select>

				<!-- Clear completed button -->
				{#if activeTab === 'done' && doneCount > 0}
					<button class="btn ml-auto btn-sm btn-error" onclick={handleClearCompleted}>
						Clear Completed
					</button>
				{/if}
			</div>

			{#if filteredTasks().length === 0}
				<p class="text-center opacity-70">No tasks found.</p>
			{:else}
				<ul
					class="flex flex-col gap-2"
					use:dndzone={{ items: dndItems, flipDurationMs: 200 }}
					onconsider={handleDndConsider}
					onfinalize={handleDndFinalize}
				>
					{#each dndItems as t (t.id)}
						<li
							class={`card border-l-4 bg-base-200 p-4 ${priorityBorder[t.Priority] || 'border-l-info'}`}
						>
							<div class="flex items-start justify-between gap-3">
								<div class="flex-1">
									<div class="flex items-center gap-2">
										<p class="font-semibold" class:line-through={t.IsDone}>{t.Title}</p>
										<span class={`badge badge-xs ${priorityBadge[t.Priority] || 'badge-info'}`}
											>{t.Priority}</span
										>
										{#if t.Category}
											<span class="badge badge-outline badge-xs">{t.Category}</span>
										{/if}
									</div>
									{#if t.Description}
										<p class="text-sm opacity-80">{t.Description}</p>
									{/if}
									<div class="mt-1 flex items-center gap-2 text-xs opacity-60">
										{#if t.DueDate}
											<span>Due: {t.DueDate}</span>
										{/if}
										{#if t.DueWarning === 'overdue'}
											<span class="badge badge-xs badge-error">Overdue</span>
										{:else if t.DueWarning === 'due_today'}
											<span class="badge badge-xs badge-warning">Due Today</span>
										{:else if t.DueWarning === 'due_soon'}
											<span
												class="badge badge-xs"
												style="background: oklch(85% 0.15 85); color: oklch(30% 0.05 85);"
												>Due Soon</span
											>
										{/if}
										{#if t.IsDone}
											<span class="text-success">Done</span>
										{/if}
									</div>

									<!-- Subtask toggle -->
									{#if t.ChildTasks && t.ChildTasks.length > 0}
										<button class="btn mt-1 btn-ghost btn-xs" onclick={() => toggleExpand(t.Id)}>
											{expandedTasks.has(t.Id) ? '▼' : '▶'}
											{t.ChildTasks.length} subtask{t.ChildTasks.length === 1 ? '' : 's'}
										</button>
									{/if}
								</div>
								<div class="flex gap-1">
									{#if !t.IsDone}
										<button
											class="btn btn-ghost btn-sm"
											title="Add Subtask"
											aria-label="Add subtask"
											onclick={() => handleAddSubtask(t.Id)}
										>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M12 4v16m8-8H4"
												/>
											</svg>
										</button>
									{/if}
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
									{#if !t.IsDone}
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
									{/if}
								</div>
							</div>

							<!-- Expanded subtasks -->
							{#if expandedTasks.has(t.Id) && t.ChildTasks && t.ChildTasks.length > 0}
								<ul class="mt-2 ml-6 flex flex-col gap-1 border-l-2 border-base-300 pl-3">
									{#each t.ChildTasks as child}
										<li
											class={`card border-l-2 bg-base-100 p-3 ${priorityBorder[child.Priority] || 'border-l-info'}`}
										>
											<div class="flex items-start justify-between gap-2">
												<div class="flex-1">
													<div class="flex items-center gap-2">
														<p class="text-sm font-medium" class:line-through={child.IsDone}>
															{child.Title}
														</p>
														<span
															class={`badge badge-xs ${priorityBadge[child.Priority] || 'badge-info'}`}
															>{child.Priority}</span
														>
													</div>
													{#if child.Description}
														<p class="text-xs opacity-80">{child.Description}</p>
													{/if}
													<div class="flex items-center gap-2 text-xs opacity-60">
														{#if child.DueDate}
															<span>Due: {child.DueDate}</span>
														{/if}
														{#if child.DueWarning === 'overdue'}
															<span class="badge badge-xs badge-error">Overdue</span>
														{:else if child.DueWarning === 'due_today'}
															<span class="badge badge-xs badge-warning">Due Today</span>
														{:else if child.DueWarning === 'due_soon'}
															<span
																class="badge badge-xs"
																style="background: oklch(85% 0.15 85); color: oklch(30% 0.05 85);"
																>Due Soon</span
															>
														{/if}
														{#if child.IsDone}
															<span class="text-success">Done</span>
														{/if}
													</div>
												</div>
												<div class="flex gap-1">
													<button
														class="btn btn-ghost btn-xs"
														title="Edit"
														aria-label="Edit subtask"
														onclick={() => handleEdit(child)}
													>
														<svg
															class="h-3 w-3"
															fill="none"
															stroke="currentColor"
															viewBox="0 0 24 24"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
															/>
														</svg>
													</button>
													<button
														class="btn btn-ghost btn-xs"
														title="Delete"
														aria-label="Delete subtask"
														onclick={() => handleDelete(child.Id)}
													>
														<svg
															class="h-3 w-3"
															fill="none"
															stroke="currentColor"
															viewBox="0 0 24 24"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
															/>
														</svg>
													</button>
													{#if !child.IsDone}
														<button
															class="btn btn-ghost btn-xs"
															title="Mark done"
															aria-label="Mark subtask done"
															onclick={() => handleMarkDone(child.Id)}
														>
															<svg
																class="h-3 w-3"
																fill="none"
																stroke="currentColor"
																viewBox="0 0 24 24"
															>
																<path
																	stroke-linecap="round"
																	stroke-linejoin="round"
																	stroke-width="2"
																	d="M5 13l4 4L19 7"
																/>
															</svg>
														</button>
													{/if}
												</div>
											</div>
										</li>
									{/each}
								</ul>
							{/if}
						</li>
					{/each}
				</ul>
			{/if}
		</div>
	</div>
</div>
