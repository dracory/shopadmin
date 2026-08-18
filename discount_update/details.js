function initDetailsApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initDetailsApp, 100);
        return;
    }

    const { createApp } = Vue;

    const app = createApp({
        data() {
            return {
                loading: false,
                discountID: typeof discountId !== 'undefined' ? discountId : '',
                code: '',
                title: '',
                description: '',
                type: '',
                amount: '',
                status: '',
                startsAt: '',
                endsAt: '',
                memo: '',
                maxUses: 0,
                maxUsesCount: 0,
                maxUsesPerCustomer: 0
            };
        },
        methods: {
            async loadDetails() {
                this.loading = true;
                try {
                    const response = await fetch(urlDetailsLoad, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ discount_id: this.discountID })
                    });

                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        if (data.status === 'success') {
                            this.code = data.data.code ?? '';
                            this.title = data.data.title ?? '';
                            this.description = data.data.description ?? '';
                            this.type = data.data.type ?? '';
                            this.amount = data.data.amount !== null && data.data.amount !== undefined ? String(data.data.amount) : '';
                            this.status = data.data.status ?? '';
                            this.startsAt = data.data.starts_at ?? '';
                            this.endsAt = data.data.ends_at ?? '';
                            this.memo = data.data.memo ?? '';
                            this.maxUses = data.data.max_uses ?? 0;
                            this.maxUsesCount = data.data.max_uses_count ?? 0;
                            this.maxUsesPerCustomer = data.data.max_uses_per_customer ?? 0;
                        } else {
                            Notiflix.Notify.failure(data.message || 'Failed to load details');
                        }
                    } else {
                        throw new Error('Invalid response format: ' + contentType);
                    }
                } catch (error) {
                    Notiflix.Notify.failure('Failed to load details: ' + error.message);
                } finally {
                    this.loading = false;
                }
            },

            async saveDetails() {
                if (!this.code) { Notiflix.Notify.failure('Code is required'); return; }
                if (!this.title) { Notiflix.Notify.failure('Title is required'); return; }
                if (!this.type) { Notiflix.Notify.failure('Type is required'); return; }
                if (!this.amount) { Notiflix.Notify.failure('Amount is required'); return; }
                if (isNaN(parseFloat(this.amount))) { Notiflix.Notify.failure('Amount must be numeric'); return; }
                if (parseFloat(this.amount) < 0) { Notiflix.Notify.failure('Amount cannot be negative'); return; }

                this.loading = true;
                try {
                    const payload = {
                        discount_id: this.discountID,
                        code: this.code,
                        title: this.title,
                        description: this.description,
                        type: this.type,
                        amount: this.amount,
                        status: this.status,
                        starts_at: this.startsAt,
                        ends_at: this.endsAt,
                        memo: this.memo,
                        max_uses: parseInt(this.maxUses) || 0,
                        max_uses_per_customer: parseInt(this.maxUsesPerCustomer) || 0
                    };
                    console.log('[discount_update] saveDetails payload:', JSON.stringify(payload, null, 2));

                    const response = await fetch(urlDetailsSave, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(payload)
                    });

                    console.log('[discount_update] saveDetails response status:', response.status, response.statusText);
                    const data = await response.json();
                    console.log('[discount_update] saveDetails response data:', JSON.stringify(data, null, 2));
                    if (data.status === 'success') {
                        Notiflix.Notify.success(data.message || 'Details saved successfully');
                    } else {
                        Notiflix.Notify.failure(data.message || 'Failed to save details');
                    }
                } catch (error) {
                    Notiflix.Notify.failure('Failed to save details: ' + error.message);
                } finally {
                    this.loading = false;
                }
            },
        },
        mounted() {
            this.loadDetails();
        }
    });

    app.use(ElementPlus);
    app.mount('#details-app');
}

document.addEventListener('DOMContentLoaded', initDetailsApp);
