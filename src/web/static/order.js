const orderIdInput = document.getElementById("order-id");
const orderOutput = document.getElementById("order");

async function fetchOrder() {
    try {
        const response = await fetch(`http://localhost:8081/order/${orderIdInput.value}`);
        const data = await response.json();
        const formattedJson = JSON.stringify(data, null, 2);

        orderOutput.innerHTML = `<pre>${formattedJson}</pre>`;
        orderOutput.style.display = "block";
    } catch (err) {
        orderOutput.textContent = "Ошибка при загрузке данных";
        console.error(err);
    }
}
