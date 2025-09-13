<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';
</script>

<svelte:head>
	<title>Towsey</title>
</svelte:head>
<PageWrapper title="Towsey" {pb}>
	<main class="grid h-full content-center justify-items-center max-lg:grid-rows-3 lg:grid-cols-3">
		<a
			href="/towelie"
			class="focus:border-primary/30 active:border-primary/30 max-lg:hover:bg-primary grid w-full content-center gap-2 rounded-xl text-center focus:border active:border"
		>
			<h2 class="font-towelie-logo hover:text-primary text-center text-7xl xl:text-9xl">Towelie</h2>
			<div>Track towel washes</div>
		</a>
		<a
			href="/nosey"
			class="focus:border-primary/30 active:border-primary/30 max-lg:hover:bg-primary grid w-full content-center gap-2 rounded-xl text-center focus:border active:border"
		>
			<h2 class="font-space hover:text-primary text-7xl font-bold xl:text-9xl">Nosey</h2>
			<div>Track nasal sprays</div>
		</a>
		<a
			href="/count"
			class="focus:border-primary/30 active:border-primary/30 max-lg:hover:bg-primary grid w-full content-center gap-2 rounded-xl text-center focus:border active:border"
		>
			<h2 class="hover:text-primary text-7xl font-black tracking-tighter xl:text-9xl">Count</h2>
			<div>Countdown for video watching</div>
		</a>
	</main>
</PageWrapper>
