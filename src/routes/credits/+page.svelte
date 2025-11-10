<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
</script>

<svelte:head>
	<title>Credits</title>
</svelte:head>

<PageWrapper title="Credits" {pb}>
	<main class="grid h-full content-center justify-items-center">
		<ul>
			<li>
				Count uses content produced with Typecast, an AI virtual actor service. Characters casted:
				Frankenstein, Neoguard from <a
					href="https://typecast.ai"
					class="hover:text-primary underline">https://typecast.ai</a
				>
			</li>
			<li>
				Alternate sounds from <a href="https://ttsmp3.com/" class="hover:text-primary underline"
					>https://ttsmp3.com/</a
				>
			</li>
		</ul>
	</main>
</PageWrapper>
