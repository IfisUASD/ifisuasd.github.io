
document.addEventListener('DOMContentLoaded', () => {
    const list = document.getElementById('publications-list');
    if (!list) return;

    const controls = document.getElementById('publications-controls');
    if (controls) controls.classList.remove('hidden');

    const items = Array.from(list.querySelectorAll('.publication-item'));
    const totalCount = document.getElementById('stats-total');
    const filteredStats = document.getElementById('stats-filtered');
    const countFiltered = document.getElementById('count-filtered');
    const paginationControls = document.getElementById('pagination-controls');

    // Filters
    const filterAuthor = document.getElementById('filter-author');
    const filterYear = document.getElementById('filter-year');
    const filterArea = document.getElementById('filter-area');

    // Sort
    const sortBy = document.getElementById('sort-by');

    // Limit
    const limitSelect = document.getElementById('limit-select');

    // State
    let currentPage = 1;
    let itemsPerPage = 10;
    let filteredItems = [...items];

    // Global data
    let jsonFilterData = { authors: [], publications: {} };

    // Initialize Filters
    initializeFilters();

    // Event Listeners
    filterAuthor.addEventListener('change', applyFilters);
    filterYear.addEventListener('change', applyFilters);
    filterArea.addEventListener('change', applyFilters);
    sortBy.addEventListener('change', () => {
        applySorting();
        render();
    });
    limitSelect.addEventListener('change', (e) => {
        itemsPerPage = e.target.value === 'all' ? Infinity : parseInt(e.target.value);
        currentPage = 1;
        render();
    });

    // Initial Render
    applyFilters();

    function initializeFilters() {
        // 1. Try to load JSON data
        const dataScript = document.getElementById('publications-data');
        if (dataScript && dataScript.textContent) {
            try {
                jsonFilterData = JSON.parse(dataScript.textContent);
            } catch (e) {
                console.error("Error parsing filter JSON:", e);
            }
        }

        const authors = jsonFilterData.authors || []; // Use JSON source strictly for authors
        const years = new Set();
        const areas = new Set();

        // Populate Years and Areas from DOM (these are reliable enough usually)
        items.forEach(item => {
             const year = item.dataset.year;
            if (year) years.add(year);

            const area = item.dataset.area;
            if (area) areas.add(area);
        });

        // Populate Author Dropdown
        authors.forEach(author => {
            const option = document.createElement('option');
            option.value = author.toLowerCase();
            option.textContent = author;
            filterAuthor.appendChild(option);
        });

        // Toggle Author Visibility
        const containerAuthor = document.getElementById('filter-container-author');
        if (containerAuthor) {
            if (authors.length === 0) containerAuthor.classList.add('hidden');
            else containerAuthor.classList.remove('hidden');
        }

        // Populate Year
        const sortedYears = Array.from(years).sort().reverse();
        sortedYears.forEach(year => {
            const option = document.createElement('option');
            option.value = year;
            option.textContent = year;
            filterYear.appendChild(option);
        });

        // Toggle Year Visibility
        const containerYear = document.getElementById('filter-container-year');
        if (containerYear) {
            if (sortedYears.length === 0) containerYear.classList.add('hidden');
            else containerYear.classList.remove('hidden');
        }

        // Populate Area
        const sortedAreas = Array.from(areas).sort();
        sortedAreas.forEach(area => {
            if (!area) return;
            const option = document.createElement('option');
            option.value = area;
            option.textContent = area;
            filterArea.appendChild(option);
        });

        // Toggle Area Visibility
        const containerArea = document.getElementById('filter-container-area');
        if (containerArea) {
            if (sortedAreas.length === 0) containerArea.classList.add('hidden');
            else containerArea.classList.remove('hidden');
        }
    }

    function applyFilters() {
        const authorValue = filterAuthor.value.toLowerCase();
        const yearValue = filterYear.value;
        const areaValue = filterArea.value;

        filteredItems = items.filter(item => {
            // Check Author Match using JSON Map if author filter is active
            let authorMatch = true;
            if (authorValue) {
                const doi = item.dataset.doi;
                // Lookup authors for this DOI in the JSON map
                const itemAuthors = (jsonFilterData.publications && jsonFilterData.publications[doi]) || [];
                // Check if any of these authors matches the selected value (case-insensitive)
                authorMatch = itemAuthors.some(a => a.toLowerCase() === authorValue);
            }

            const yearMatch = !yearValue || item.dataset.year === yearValue;
            const areaMatch = !areaValue || item.dataset.area === areaValue;
            return authorMatch && yearMatch && areaMatch;
        });

        // Update Stats
        if (filteredItems.length !== items.length) {
            filteredStats.classList.remove('hidden');
            countFiltered.textContent = filteredItems.length;
        } else {
            filteredStats.classList.add('hidden');
        }

        currentPage = 1;
        applySorting();
        render();
    }

    function applySorting() {
        const sortValue = sortBy.value;
        filteredItems.sort((a, b) => {
            const dateA = new Date(a.dataset.date);
            const dateB = new Date(b.dataset.date);
            const citA = parseInt(a.dataset.citations) || 0;
            const citB = parseInt(b.dataset.citations) || 0;
            const titleA = a.dataset.title.toLowerCase();
            const titleB = b.dataset.title.toLowerCase();

            switch (sortValue) {
                case 'date-asc': return dateA - dateB;
                case 'date-desc': return dateB - dateA;
                case 'citations-asc': return citA - citB;
                case 'citations-desc': return citB - citA;
                case 'title-asc': return titleA.localeCompare(titleB);
                case 'title-desc': return titleB.localeCompare(titleA);
                default: return 0;
            }
        });
    }

    function render() {
        // Clear list
        list.innerHTML = '';

        // Pagination
        const start = (currentPage - 1) * itemsPerPage;
        const end = itemsPerPage === Infinity ? filteredItems.length : start + itemsPerPage;
        const visibleItems = filteredItems.slice(start, end);

        visibleItems.forEach(item => list.appendChild(item));

        renderPaginationControls();
    }

    function renderPaginationControls() {
        paginationControls.innerHTML = '';
        if (itemsPerPage === Infinity || filteredItems.length <= itemsPerPage) {
            paginationControls.classList.add('hidden');
            return;
        }

        paginationControls.classList.remove('hidden');
        const totalPages = Math.ceil(filteredItems.length / itemsPerPage);

        // Previous
        const prevBtn = document.createElement('button');
        prevBtn.className = 'join-item btn btn-sm';
        prevBtn.textContent = '«';
        prevBtn.disabled = currentPage === 1;
        prevBtn.onclick = () => { currentPage--; render(); };
        paginationControls.appendChild(prevBtn);

        // Page Numbers
        let startPage = Math.max(1, currentPage - 2);
        let endPage = Math.min(totalPages, startPage + 4);
        if (endPage - startPage < 4) {
             startPage = Math.max(1, endPage - 4);
        }

        for (let i = startPage; i <= endPage; i++) {
            const btn = document.createElement('button');
            btn.className = `join-item btn btn-sm ${i === currentPage ? 'btn-active' : ''}`;
            btn.textContent = i;
            btn.onclick = () => { currentPage = i; render(); };
            paginationControls.appendChild(btn);
        }

        // Next
        const nextBtn = document.createElement('button');
        nextBtn.className = 'join-item btn btn-sm';
        nextBtn.textContent = '»';
        nextBtn.disabled = currentPage === totalPages;
        nextBtn.onclick = () => { currentPage++; render(); };
        paginationControls.appendChild(nextBtn);
    }
});
