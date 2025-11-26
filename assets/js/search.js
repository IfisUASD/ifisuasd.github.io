document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.getElementById('searchInput');
    const searchResults = document.getElementById('searchResults');
    const loading = document.getElementById('searchLoading');
    const noResults = document.getElementById('searchNoResults');
    
    let searchIndex = [];
    let isLoaded = false;

    // Determine current language from URL
    const isEnglish = window.location.pathname.startsWith('/en');
    const indexUrl = isEnglish ? '/en/search.json' : '/search.json';

    // Load Index
    async function loadIndex() {
        try {
            loading.classList.remove('hidden');
            const response = await fetch(indexUrl);
            searchIndex = await response.json();
            isLoaded = true;
            loading.classList.add('hidden');
            // Trigger search if user already typed
            if (searchInput.value.trim()) {
                performSearch(searchInput.value);
            }
        } catch (error) {
            console.error('Error loading search index:', error);
            loading.classList.add('hidden');
        }
    }

    // Perform Search
    function performSearch(query) {
        if (!query.trim()) {
            searchResults.innerHTML = '';
            noResults.classList.add('hidden');
            return;
        }

        const lowerQuery = query.toLowerCase();
        const results = searchIndex.filter(item => {
            return item.title.toLowerCase().includes(lowerQuery) || 
                   item.summary.toLowerCase().includes(lowerQuery) ||
                   item.type.toLowerCase().includes(lowerQuery) ||
                   (item.tags && item.tags.some(tag => tag.toLowerCase().includes(lowerQuery)));
        });

        renderResults(results);
    }

    // Render Results
    function renderResults(results) {
        searchResults.innerHTML = '';
        
        if (results.length === 0) {
            noResults.classList.remove('hidden');
            return;
        }
        
        noResults.classList.add('hidden');

        results.forEach(item => {
            const div = document.createElement('div');
            div.className = 'card bg-base-100 shadow-sm border border-base-200 hover:border-primary transition-colors';
            div.innerHTML = `
                <div class="card-body p-4">
                    <div class="flex justify-between items-start">
                        <h3 class="card-title text-lg">
                            <a href="${item.url}" class="hover:text-primary hover:underline">${item.title}</a>
                        </h3>
                        <span class="badge badge-ghost text-xs">${item.type}</span>
                    </div>
                    <p class="text-sm text-base-content/70 line-clamp-2">${item.summary}</p>
                </div>
            `;
            searchResults.appendChild(div);
        });
    }

    // Event Listeners
    searchInput.addEventListener('input', (e) => {
        if (!isLoaded) {
            loadIndex();
        } else {
            performSearch(e.target.value);
        }
    });

    // Load index on focus if not loaded
    searchInput.addEventListener('focus', () => {
        if (!isLoaded) {
            loadIndex();
        }
    });
});
