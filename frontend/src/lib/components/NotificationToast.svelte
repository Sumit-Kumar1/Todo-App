<script lang="ts">
	import { notifications, type ToastNotification } from '$lib/stores/notifications';

	function bg(level: ToastNotification['level']): string {
		switch (level) {
			case 'success':
				return 'bg-emerald-500 text-white';
			case 'warning':
				return 'bg-amber-500 text-black';
			case 'error':
				return 'bg-rose-500 text-white';
			default:
				return 'bg-slate-800 text-white';
		}
	}
</script>

<div class="pointer-events-none fixed inset-0 z-50">
	<div class="absolute top-4 right-4 flex flex-col items-end gap-3">
		{#each $notifications as n (n.id)}
			<div
				class={`pointer-events-auto rounded-2xl shadow-2xl/60 ring-1 shadow-black/40 ring-black/10 ${bg(n.level)} w-[360px] max-w-[90vw] overflow-hidden`}
				role="status"
				aria-live="polite"
			>
				<div class="p-4 text-sm leading-5">{n.message}</div>
			</div>
		{/each}
	</div>
</div>

<style>
	/* simple slide-in */
	:global(.toast-enter) {
		transform: translateY(-8px);
		opacity: 0;
	}
	:global(.toast-enter-active) {
		transition: all 180ms ease;
		transform: translateY(0);
		opacity: 1;
	}
</style>
