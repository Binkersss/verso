let manifest = {};
let currentRoute = 'about';

async function init() {
	try {
		const response = await fetch('/content.json');
		manifest = await response.json();

		renderNav();

		const path = window.location.pathname.slice(1) || 'about';
		navigateTo(path, false);

		window.addEventListener('popstate', () => {
			const path = window.location.pathname.slice(1) || 'about';
			navigateTo(path, false);
		});
	} catch (error) {
		console.error('Failed to load content:', error);
		document.getElementById('content').innerHTML = '<p>Failed to load content</p>';
	}
}

function renderNav() {
	const nav = document.getElementById('nav');
	nav.innerHTML = '';

	Object.keys(manifest.pages).forEach(route => {
		const item = document.createElement('div');
		item.className = 'nav-item';
		item.textContent = route;
		item.onclick = () => navigateTo(route);
		nav.appendChild(item);
	});
}

function navigateTo(route, pushState = true) {
	const page = manifest.pages[route];
	if (!page) {
		console.error('Page not found:', route);
		return;
	}

	currentRoute = route;

	const content = document.getElementById('content');
	content.classList.add('fade');

	setTimeout(() => {
		content.innerHTML = page.content;
		content.classList.remove('fade');

		document.querySelectorAll('.nav-item').forEach(item => {
			item.classList.toggle('active', item.textContent === route);
		});

		if (pushState) {
			window.history.pushState({}, '', `/${route}`);
		}

		window.scrollTo(0, 0);
	}, 150);
}

init();
