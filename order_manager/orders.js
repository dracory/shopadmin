const { createApp } = Vue;

const OrdersApp = {
    data() {
        return {
            // UI state
            loading: true,
            showFilterModal: false,

            // Order data
            orders: [],
            totalOrders: 0,

            // Pagination
            currentPage: 0,
            perPage: 10,
            jumpToPage: 1,

            // Applied conditions (sent to the server)
            conditions: [],
            // Modal state — cloned from conditions when opening
            modalConditions: [],

            // Sorting
            sortByColumn: 'created_at',
            sortOrder: 'desc',

            // Filter condition options — must match the Go handler fields.
            conditionOptions: [
                { value: 'order_id', label: 'Order ID', operators: ['contains'], inputType: 'text', placeholder: 'order id...' },
                { value: 'customer_name', label: 'Customer Name', operators: ['contains'], inputType: 'text', placeholder: 'customer name...' },
                { value: 'customer_email', label: 'Customer Email', operators: ['contains'], inputType: 'text', placeholder: 'customer email...' },
                { value: 'status', label: 'Status', operators: ['equals'], inputType: 'select', placeholder: '', options: [
                    { value: 'completed', label: 'Completed' },
                    { value: 'pending', label: 'Pending' },
                    { value: 'cancelled', label: 'Cancelled' },
                ] },
                { value: 'date_from', label: 'Created after', operators: ['equals'], inputType: 'date', placeholder: '' },
                { value: 'date_to', label: 'Created before', operators: ['equals'], inputType: 'date', placeholder: '' },
            ],
            opLabels: { equals: '=', contains: 'contains' }
        };
    },

    computed: {
        /**
         * Returns total number of pages.
         */
        totalPages() {
            return Math.ceil(this.totalOrders / this.perPage);
        },

        /**
         * Returns a human-readable string describing the current filter state.
         */
        filterStatus() {
            if (this.conditions.length === 0) return 'Showing all orders';
            return this.conditions.length + ' filter(s) active';
        },

        /**
         * Returns true if any filters are active.
         */
        hasActiveFilters() {
            return this.conditions.length > 0;
        }
    },

    mounted() {
        // Read filters from URL parameters for shareable URLs
        const urlParams = new URLSearchParams(window.location.search);

        // Try JSON filters first (new format)
        const filtersRaw = urlParams.get('filters');
        if (filtersRaw) {
            try {
                const parsed = JSON.parse(filtersRaw);
                if (Array.isArray(parsed)) {
                    this.conditions = parsed.filter(c => c && c.field && c.value);
                }
            } catch (e) { /* ignore malformed */ }
        } else {
            // Fallback: legacy individual params
            if (urlParams.get('status')) this.conditions.push({ field: 'status', operator: 'equals', value: urlParams.get('status') });
            if (urlParams.get('customer_name')) this.conditions.push({ field: 'customer_name', operator: 'contains', value: urlParams.get('customer_name') });
            if (urlParams.get('customer_email')) this.conditions.push({ field: 'customer_email', operator: 'contains', value: urlParams.get('customer_email') });
            if (urlParams.get('order_id')) this.conditions.push({ field: 'order_id', operator: 'contains', value: urlParams.get('order_id') });
            if (urlParams.get('created_from')) this.conditions.push({ field: 'date_from', operator: 'equals', value: urlParams.get('created_from') });
            if (urlParams.get('created_to')) this.conditions.push({ field: 'date_to', operator: 'equals', value: urlParams.get('created_to') });
        }

        this.sortByColumn = urlParams.get('sort_by') || 'created_at';
        this.sortOrder = urlParams.get('sort_order') || 'desc';
        this.currentPage = parseInt(urlParams.get('page') || '0', 10);
        this.perPage = parseInt(urlParams.get('per_page') || '10', 10);

        this.loadOrders();
    },
    methods: {
        /**
         * Loads orders from the server based on current filters, pagination, and sorting.
         */
        async loadOrders() {
            this.loading = true;
            try {
                const response = await fetch(urlLoadOrders, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        page: this.currentPage,
                        per_page: this.perPage,
                        conditions: this.conditions,
                        sort_by: this.sortByColumn,
                        sort_order: this.sortOrder
                    })
                });
                const data = await response.json();

                if (data.status === 'success') {
                    this.orders = data.data?.orders || [];
                    this.totalOrders = data.data?.total || 0;
                    // Sync jumpToPage with currentPage to ensure UI consistency
                    this.jumpToPage = this.currentPage + 1;
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to load orders'
                    });
                }
            } catch (error) {
                console.error('Error loading orders:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'Failed to load orders'
                });
            } finally {
                this.loading = false;
            }
        },
        /**
         * Sorts the orders by the specified column.
         */
        sortBy(column) {
            if (this.sortByColumn === column) {
                this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
            } else {
                this.sortByColumn = column;
                this.sortOrder = 'asc';
            }
            
            // For computed fields (customer_name, customer_email), sort client-side
            if (column === 'customer_name' || column === 'customer_email') {
                this.sortOrdersClientSide();
            } else {
                this.currentPage = 0;
                this.applyFilters();
            }
        },

        /**
         * Sorts orders client-side for computed fields.
         */
        sortOrdersClientSide() {
            const column = this.sortByColumn;
            const order = this.sortOrder === 'asc' ? 1 : -1;
            
            this.orders.sort((a, b) => {
                const aVal = (a[column] || '').toLowerCase();
                const bVal = (b[column] || '').toLowerCase();
                
                if (aVal < bVal) return -1 * order;
                if (aVal > bVal) return 1 * order;
                return 0;
            });
        },

        /**
         * Navigates to the specified page number.
         */
        goToPage(page) {
            if (page < 0) return;
            const maxPage = Math.ceil(this.totalOrders / this.perPage) - 1;
            if (page > maxPage) page = maxPage;
            this.currentPage = page;
            this.jumpToPage = page + 1;
            this.applyFilters();
        },

        /**
         * Changes the per-page setting and reloads orders.
         */
        changePerPage() {
            this.perPage = parseInt(this.perPage, 10);
            this.currentPage = 0;
            this.jumpToPage = 1;
            this.applyFilters();
        },

        /**
         * Opens the filter modal.
         */
        openFilterModal() {
            this.modalConditions = this.conditions.map(c => ({ ...c }));
            if (this.modalConditions.length === 0) {
                this.modalConditions.push({ field: '', operator: 'equals', value: '' });
            }
            this.showFilterModal = true;
        },

        /**
         * Closes the filter modal.
         */
        closeFilterModal() {
            this.showFilterModal = false;
        },

        addModalCondition() {
            this.modalConditions.push({ field: '', operator: 'equals', value: '' });
        },

        removeModalCondition(i) {
            this.modalConditions.splice(i, 1);
        },

        clearModalConditions() {
            this.modalConditions = [];
        },

        onFieldChange(c) {
            const opt = this.getFieldOpt(c.field);
            if (opt) {
                c.operator = opt.operators[0];
            }
            c.value = '';
        },

        applyConditions() {
            this.conditions = this.modalConditions.filter(c => c.field && String(c.value).trim());
            this.showFilterModal = false;
            this.currentPage = 0;
            this.applyFilters();
        },

        /**
         * Applies the current filters and updates the URL.
         */
        applyFilters() {
            const params = new URLSearchParams();
            const existing = new URLSearchParams(window.location.search);
            const controller = existing.get('controller');
            if (controller) params.set('controller', controller);

            if (this.conditions.length > 0) params.set('filters', JSON.stringify(this.conditions));
            params.set('page', this.currentPage);
            params.set('per_page', this.perPage);
            params.set('sort_order', this.sortOrder);
            params.set('sort_by', this.sortByColumn);

            const newUrl = `${window.location.pathname}?${params.toString()}`;
            window.history.pushState({}, '', newUrl);

            this.loadOrders();
        },

        /**
         * Clears all filters and resets to default state.
         */
        clearFilters() {
            this.conditions = [];
            this.currentPage = 0;
            this.applyFilters();
        },

        // === Filter option helpers ===
        getFieldOpt(v) { return this.conditionOptions.find(o => o.value === v); },
        opLabel(op) { return this.opLabels[op] || op; },
        getFieldOperators(v) { const o = this.getFieldOpt(v); return o ? o.operators : ['equals']; },
        getFieldInputType(v) { const o = this.getFieldOpt(v); return o ? o.inputType : 'text'; },
        getFieldPlaceholder(v) { const o = this.getFieldOpt(v); return o ? o.placeholder : ''; },
        getFieldSelectOptions(v) { const o = this.getFieldOpt(v); return o && o.options ? o.options : []; },

        getStatusBadgeClass(status) {
            switch (status) {
                case 'completed':
                    return 'bg-success';
                case 'pending':
                    return 'bg-warning';
                case 'cancelled':
                    return 'bg-danger';
                default:
                    return 'bg-secondary';
            }
        },
        formatDate(dateString) {
            if (!dateString) return '';
            const date = new Date(dateString);
            
            const day = date.getDate().toString().padStart(2, '0');
            const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
            const month = months[date.getMonth()];
            const year = date.getFullYear();
            
            const hours = date.getHours().toString().padStart(2, '0');
            const minutes = date.getMinutes().toString().padStart(2, '0');
            
            return `${day} ${month} ${year}<br><span class="small text-muted">${hours}:${minutes}</span>`;
        }
    }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    const el = document.getElementById('app');
    if (el) {
        createApp(OrdersApp).mount('#app');
    }
});
