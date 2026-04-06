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
			if (res.ok) {
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
					t.Title.toLowerCase().includes(search) ||
					(t.Description && t.Description.toLowerCase().includes(search))
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
	let activeCount = $derived($tasks.filter((t) => !t.ParentId && !t.IsDone).length);
	let totalCount = $derived($tasks.filter((t) => !t.ParentId).length);

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
		LOW: 'border-l-base-content/20',
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

<!-- Navbar -->
<nav class="navbar sticky top-0 z-30 border-b border-base-300 bg-base-100/80 px-4 backdrop-blur-md">
	<div class="flex-1 gap-2">
		<span class="text-lg font-bold sm:text-xl">Todo App</span>
	</div>
	<div class="flex-none gap-1">
		<button class="btn btn-circle btn-ghost btn-sm" aria-label="toggle theme" onclick={toggleTheme}>
			{#if $theme === 'light'}
				<Fa icon={faMoon} />
			{:else}
				<Fa icon={faSun} />
			{/if}
		</button>
		<button class="btn btn-ghost btn-sm" aria-label="logout" onclick={userLogout}>
			<Fa icon={faArrowRightFromBracket}></Fa>
			<span class="hidden sm:inline">Logout</span>
		</button>
	</div>
</nav>

<main class="mx-auto w-full max-w-2xl px-4 py-6">
	<!-- Header with stats and create button -->
	<div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h1 class="text-2xl font-bold">My Tasks</h1>
			<p class="text-sm text-base-content/60">
				{activeCount} active, {doneCount} completed
			</p>
		</div>
		<button class="btn btn-primary sm:btn-md" aria-label="task-create" onclick={openCreateModal}>
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			New Task
		</button>
	</div>

	<TodoForms bind:this={todoForm} />

	<!-- Search bar -->
	<div class="mb-4">
		<label class="input-bordered input flex w-full items-center gap-2">
			<svg class="h-4 w-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
				/>
			</svg>
			<input
				type="text"
				class="grow"
				placeholder="Search tasks..."
				value={searchQuery}
				oninput={onSearchInput}
			/>
		</label>
	</div>

	<!-- Filter row -->
	<div class="mb-4 flex flex-wrap items-center gap-2">
		<!-- Tab filters -->
		<div role="tablist" class="tabs-box tabs tabs-sm">
			<button
				role="tab"
				class={`tab ${activeTab === 'all' ? 'tab-active' : ''}`}
				onclick={() => (activeTab = 'all')}>All ({totalCount})</button
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

		<div class="flex flex-1 flex-wrap items-center justify-end gap-2">
			<!-- Priority filter -->
			<select class="select-bordered select select-sm" bind:value={filterPriority}>
				<option value="">Priority</option>
				<option value="URGENT">Urgent</option>
				<option value="HIGH">High</option>
				<option value="MEDIUM">Medium</option>
				<option value="LOW">Low</option>
			</select>

			<!-- Category filter -->
			<select class="select-bordered select select-sm" bind:value={filterCategory}>
				<option value="">Category</option>
				{#each allCategories() as cat}
					<option value={cat}>{cat}</option>
				{/each}
			</select>

			<!-- Clear completed button -->
			{#if activeTab === 'done' && doneCount > 0}
				<button class="btn btn-outline btn-sm btn-error" onclick={handleClearCompleted}>
					Clear Done
				</button>
			{/if}
		</div>
	</div>

	<!-- Task list -->
	{#if filteredTasks().length === 0}
		<div class="flex flex-col items-center justify-center py-16 text-base-content/40">
			<svg class="mb-4 h-16 w-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="1.5"
					d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
				/>
			</svg>
			<p class="text-lg font-medium">No tasks found</p>
			<p class="text-sm">
				{#if activeTab === 'done'}
					No completed tasks yet
				{:else if debouncedSearch}
					Try a different search term
				{:else}
					Create your first task to get started
				{/if}
			</p>
		</div>
	{:else}
		<ul
			class="flex flex-col gap-3"
			use:dndzone={{ items: dndItems, flipDurationMs: 200 }}
			onconsider={handleDndConsider}
			onfinalize={handleDndFinalize}
		>
			{#each dndItems as t (t.id)}
				<li
					class={`card border-l-4 bg-base-200 transition-shadow hover:shadow-md ${priorityBorder[t.Priority] || 'border-l-info'}`}
				>
					<div class="card-body gap-2 p-4">
						<!-- Top row: title + badges -->
						<div class="flex flex-wrap items-start justify-between gap-2">
							<div class="min-w-0 flex-1">
								<div class="flex flex-wrap items-center gap-2">
									<h3
										class="truncate text-base font-semibold"
										class:line-through={t.IsDone}
										class:opacity-50={t.IsDone}
									>
										{t.Title}
									</h3>
									<span class={`badge badge-sm ${priorityBadge[t.Priority] || 'badge-info'}`}
										>{t.Priority}</span
									>
									{#if t.Category}
										<span class="badge badge-outline badge-sm">{t.Category}</span>
									{/if}
								</div>
								{#if t.Description}
									<p class="mt-1 text-sm text-base-content/70">{t.Description}</p>
								{/if}
							</div>

							<!-- Actions dropdown for mobile, inline for desktop -->
							<div class="dropdown dropdown-end sm:hidden">
								<div tabindex="0" role="button" class="btn btn-circle btn-ghost btn-sm">
									<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M12 5v.01M12 12v.01M12 19v.01"
										/>
									</svg>
								</div>
								<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
								<ul
									tabindex="0"
									class="dropdown-content menu z-10 w-44 rounded-box bg-base-100 p-2 shadow-lg"
								>
									{#if !t.IsDone}
										<li><button onclick={() => handleAddSubtask(t.Id)}>Add Subtask</button></li>
									{/if}
									<li><button onclick={() => handleEdit(t)}>Edit</button></li>
									{#if !t.IsDone}
										<li><button onclick={() => handleMarkDone(t.Id)}>Mark Done</button></li>
									{/if}
									<li>
										<button class="text-error" onclick={() => handleDelete(t.Id)}>Delete</button>
									</li>
								</ul>
							</div>

							<!-- Desktop action buttons -->
							<div class="hidden gap-1 sm:flex">
								{#if !t.IsDone}
									<button
										class="btn btn-circle btn-ghost btn-sm"
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
									class="btn btn-circle btn-ghost btn-sm"
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
									class="btn btn-circle btn-ghost btn-sm"
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
										class="btn btn-circle btn-ghost btn-sm"
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

						<!-- Meta row: due date + warnings -->
						<div class="flex flex-wrap items-center gap-2 text-xs text-base-content/50">
							{#if t.DueDate}
								<span>Due: {t.DueDate}</span>
							{/if}
							{#if t.DueWarning === 'overdue'}
								<span class="badge badge-sm badge-error">Overdue</span>
							{:else if t.DueWarning === 'due_today'}
								<span class="badge badge-sm badge-warning">Due Today</span>
							{:else if t.DueWarning === 'due_soon'}
								<span
									class="badge badge-sm"
									style="background: oklch(85% 0.15 85); color: oklch(30% 0.05 85);">Due Soon</span
								>
							{/if}
							{#if t.IsDone}
								<span class="badge badge-sm badge-success">Done</span>
							{/if}
						</div>

						<!-- Subtask toggle -->
						{#if t.ChildTasks && t.ChildTasks.length > 0}
							<button
								class="btn mt-1 w-fit gap-1 btn-ghost btn-xs"
								onclick={() => toggleExpand(t.Id)}
							>
								<svg
									class="h-3 w-3 transition-transform"
									class:rotate-90={expandedTasks.has(t.Id)}
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M9 5l7 7-7 7"
									/>
								</svg>
								{t.ChildTasks.length} subtask{t.ChildTasks.length === 1 ? '' : 's'}
							</button>
						{/if}

						<!-- Expanded subtasks -->
						{#if expandedTasks.has(t.Id) && t.ChildTasks && t.ChildTasks.length > 0}
							<ul class="mt-2 flex flex-col gap-2 border-l-2 border-base-300 pl-3 sm:ml-4">
								{#each t.ChildTasks as child}
									<li
										class={`card border-l-2 bg-base-100 p-3 ${priorityBorder[child.Priority] || 'border-l-info'}`}
									>
										<div class="flex flex-wrap items-start justify-between gap-2">
											<div class="min-w-0 flex-1">
												<div class="flex flex-wrap items-center gap-2">
													<p
														class="text-sm font-medium"
														class:line-through={child.IsDone}
														class:opacity-50={child.IsDone}
													>
														{child.Title}
													</p>
													<span
														class={`badge badge-xs ${priorityBadge[child.Priority] || 'badge-info'}`}
														>{child.Priority}</span
													>
												</div>
												{#if child.Description}
													<p class="mt-0.5 text-xs text-base-content/70">{child.Description}</p>
												{/if}
												<div
													class="mt-1 flex flex-wrap items-center gap-2 text-xs text-base-content/50"
												>
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
													class="btn btn-circle btn-ghost btn-xs"
													title="Edit"
													aria-label="Edit subtask"
													onclick={() => handleEdit(child)}
												>
													<svg
														class="h-3.5 w-3.5"
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
													class="btn btn-circle btn-ghost btn-xs"
													title="Delete"
													aria-label="Delete subtask"
													onclick={() => handleDelete(child.Id)}
												>
													<svg
														class="h-3.5 w-3.5"
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
														class="btn btn-circle btn-ghost btn-xs"
														title="Mark done"
														aria-label="Mark subtask done"
														onclick={() => handleMarkDone(child.Id)}
													>
														<svg
															class="h-3.5 w-3.5"
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
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</main>
