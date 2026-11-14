function initVariantsApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initVariantsApp, 100);
        return;
    }

    const { createApp } = Vue;

    createApp({
        data() {
            return {
                loading: false,
                collapsed: false,
                dimensions: [],
                variants: {},
                showDimModal: false,
                showOptionModal: false,
                optionModalDimIndex: -1,
                newDimName: '',
                newDimOptions: '',
                newOptionValue: '',
                productID: variantsProductId
            };
        },
        computed: {
            combos() {
                // Cartesian product of all dimension options
                if (this.dimensions.length === 0) return [];
                let combos = this.dimensions[0].options.map(o => ({ [this.dimensions[0].name]: o }));
                for (let i = 1; i < this.dimensions.length; i++) {
                    const next = [];
                    combos.forEach(c => this.dimensions[i].options.forEach(o =>
                        next.push({ ...c, [this.dimensions[i].name]: o })));
                    combos = next;
                }
                return combos;
            },
            optionModalDim() {
                return this.optionModalDimIndex >= 0 && this.dimensions[this.optionModalDimIndex]
                    ? this.dimensions[this.optionModalDimIndex].name
                    : '';
            }
        },
        methods: {
            comboKey(combo) {
                // Stable key in dimension order
                return this.dimensions
                    .filter(d => combo[d.name] !== undefined)
                    .map(d => d.name + ':' + combo[d.name])
                    .join('|');
            },

            variantFor(combo) {
                const key = this.comboKey(combo);
                if (!this.variants[key]) {
                    this.variants[key] = { sku: '', price: '', quantity: '' };
                }
                return this.variants[key];
            },

            pruneVariants() {
                const valid = new Set(this.combos.map(c => this.comboKey(c)));
                Object.keys(this.variants).forEach(k => {
                    if (!valid.has(k)) delete this.variants[k];
                });
            },

            addDimension() {
                const name = this.newDimName.trim().toLowerCase();
                const options = this.newDimOptions.split(',')
                    .map(s => s.trim()).filter(Boolean);
                if (!name || options.length === 0) return;
                if (this.dimensions.some(d => d.name === name)) {
                    Notiflix.Notify.failure('Dimension "' + name + '" already exists');
                    return;
                }
                this.dimensions.push({ name, required: true, options });
                this.newDimName = '';
                this.newDimOptions = '';
                this.showDimModal = false;
            },

            removeDimension(di) {
                const dim = this.dimensions[di];
                const n = this.combos.filter(c => c[dim.name] !== undefined).length;
                Swal.fire({
                    title: 'Remove Dimension',
                    html: 'Remove dimension <b>' + dim.name + '</b> and all its options?<br>' +
                        'This deletes <b>' + n + '</b> variant row(s) — their SKU, price ' +
                        'and stock values will be lost on save.',
                    icon: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#d33',
                    confirmButtonText: 'Yes, delete it'
                }).then(result => {
                    if (!result.isConfirmed) return;
                    this.dimensions.splice(di, 1);
                    this.pruneVariants();
                });
            },

            openOptionModal(di) {
                this.optionModalDimIndex = di;
                this.newOptionValue = '';
                this.showOptionModal = true;
            },

            addOption() {
                const v = this.newOptionValue.trim();
                const dim = this.dimensions[this.optionModalDimIndex];
                if (!v || !dim) return;
                if (dim.options.includes(v)) {
                    Notiflix.Notify.failure('Option "' + v + '" already exists');
                    return;
                }
                dim.options.push(v);
                this.showOptionModal = false;
            },

            removeOption(di, oi) {
                const dim = this.dimensions[di];
                const opt = dim.options[oi];
                const n = this.combos.filter(c => c[dim.name] === opt).length;
                Swal.fire({
                    title: 'Remove Option',
                    html: 'Remove option <b>' + opt + '</b> from <b>' + dim.name + '</b>?<br>' +
                        'This deletes <b>' + n + '</b> variant row(s) — their SKU, price ' +
                        'and stock values will be lost on save.',
                    icon: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#d33',
                    confirmButtonText: 'Yes, delete it'
                }).then(result => {
                    if (!result.isConfirmed) return;
                    dim.options.splice(oi, 1);
                    if (dim.options.length === 0) this.dimensions.splice(di, 1);
                    this.pruneVariants();
                });
            },

            removeCombination(combo) {
                // Removing a row means removing the option values that
                // produced it — which also removes sibling rows.
                const removed = new Set();
                Object.entries(combo).forEach(([dn, ov]) => {
                    this.combos.forEach(c => {
                        if (c[dn] === ov) removed.add(this.comboKey(c));
                    });
                });
                Swal.fire({
                    title: 'Remove Combination',
                    html: 'Remove this combination?<br>' +
                        'This deletes its option(s) from the matrix, removing <b>' + removed.size +
                        '</b> variant row(s) total — SKU, price and stock values will be lost on save.',
                    icon: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#d33',
                    confirmButtonText: 'Yes, delete it'
                }).then(result => {
                    if (!result.isConfirmed) return;
                    Object.entries(combo).forEach(([dn, ov]) => {
                        const d = this.dimensions.find(x => x.name === dn);
                        if (d) d.options = d.options.filter(o => o !== ov);
                    });
                    this.dimensions = this.dimensions.filter(d => d.options.length > 0);
                    this.pruneVariants();
                });
            },

            async loadVariants() {
                this.loading = true;
                try {
                    const response = await fetch(variantsUrlLoad, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ product_id: this.productID })
                    });

                    const contentType = response.headers.get('content-type');
                    if (!contentType || !contentType.includes('application/json')) {
                        const text = await response.text();
                        console.error('Non-JSON response:', text);
                        throw new Error('Invalid response format: ' + contentType);
                    }

                    const data = await response.json();
                    this.dimensions = data.dimensions || [];
                    const map = {};
                    (data.variants || []).forEach(v => {
                        const key = this.dimensions
                            .filter(d => (v.values || {})[d.name] !== undefined)
                            .map(d => d.name + ':' + v.values[d.name])
                            .join('|');
                        if (key) {
                            map[key] = {
                                sku: v.sku || '',
                                price: v.price || '',
                                quantity: v.quantity || ''
                            };
                        }
                    });
                    this.variants = map;
                } catch (error) {
                    console.error('Error loading variants:', error);
                    Notiflix.Notify.failure('Failed to load variants: ' + error.message);
                } finally {
                    this.loading = false;
                }
            },

            async saveVariants() {
                this.loading = true;
                try {
                    const payload = this.combos.map(combo => {
                        const v = this.variants[this.comboKey(combo)] || {};
                        return {
                            values: combo,
                            sku: v.sku || '',
                            price: String(v.price ?? ''),
                            quantity: String(v.quantity ?? '')
                        };
                    });

                    const response = await fetch(variantsUrlSave, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            product_id: this.productID,
                            dimensions: this.dimensions,
                            variants: payload
                        })
                    });

                    const data = await response.json();
                    if (data.status === 'success') {
                        Notiflix.Notify.success(data.message || 'Variants saved successfully');
                    } else {
                        Notiflix.Notify.failure(data.message || 'Failed to save variants');
                    }
                } catch (error) {
                    console.error('Error saving variants:', error);
                    Notiflix.Notify.failure('Failed to save variants: ' + error.message);
                } finally {
                    this.loading = false;
                }
            }
        },
        mounted() {
            this.loadVariants();
        }
    }).mount('#variants-app');
}

document.addEventListener('DOMContentLoaded', initVariantsApp);
