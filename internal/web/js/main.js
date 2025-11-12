const spinner = document.getElementById("spinner");
const result = document.getElementById("result");

function fetchOrder() {
    const id = document.getElementById("orderID").value.trim();
    if (!id) return alert("Enter order UID");

    spinner.style.display = "block";
    result.innerHTML = "";

    fetch(`/order/${id}`)
        .then(res => {
            spinner.style.display = "none";
            if (!res.ok) throw new Error("Order not found");
            return res.json();
        })
        .then(order => renderOrder(order))
        .catch(err => {
            result.innerHTML = `<div class="card" style="color:red">${err.message}</div>`;
        });
}

function renderOrder(order) {
    const d = order.Delivery || order.delivery || {};
    const p = order.Payment || order.payment || {};
    const items = order.Items || order.items || [];

    const itemsHTML = items.length
        ? items.map(i => `
            <div class="item">
                <span class="label">Name:</span> <span class="value">${i.Name || i.name || "Item"}</span><br>
                <span class="label">Count:</span> <span class="value">${i.Sale || i.sale || 0}</span><br>
                <span class="label">Price:</span> <span class="value">${i.TotalPrice || i.total_price || "-"} ${p.Currency || p.currency || ""}</span>
            </div>
        `).join("")
        : `<div class="item">No items</div>`;

    result.innerHTML = `
        <div class="card">
            <div class="card-title">Order</div>
            <div><span class="label">UID:</span> <span class="value">${order.OrderUID || order.order_uid || "-"}</span></div>
            <div><span class="label">Track:</span> <span class="value">${order.TrackNumber || order.track_number || "-"}</span></div>
        </div>

        <div class="card">
            <div class="card-title">Delivery</div>
            <div><span class="label">Name:</span> <span class="value">${d.Name || d.name || ""}</span></div>
            <div><span class="label">Address:</span> <span class="value">${d.Address || d.address || ""}</span></div>
            <div><span class="label">City:</span> <span class="value">${d.City || d.city || ""}</span></div>
            <div><span class="label">Region:</span> <span class="value">${d.Region || d.region || ""}</span></div>
            <div><span class="label">Zip:</span> <span class="value">${d.Zip || d.zip || ""}</span></div>
            <div><span class="label">Phone:</span> <span class="value">${d.Phone || d.phone || ""}</span></div>
            <div><span class="label">Email:</span> <span class="value">${d.Email || d.email || ""}</span></div>
        </div>

        <div class="card">
            <div class="card-title">Payment</div>
            <div><span class="label">Provider:</span> <span class="value">${p.Provider || p.provider || ""}</span></div>
            <div><span class="label">Amount:</span> <span class="value">${p.Amount || p.amount || ""} ${p.Currency || p.currency || ""}</span></div>
            <div><span class="label">Bank:</span> <span class="value">${p.Bank || p.bank || ""}</span></div>
            <div><span class="label">Transaction:</span> <span class="value">${p.Transaction || p.transaction || ""}</span></div>
        </div>

        <div class="card">
            <div class="card-title">Items</div>
            ${itemsHTML}
        </div>
    `;
}