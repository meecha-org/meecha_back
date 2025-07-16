const popup_distance_area = document.getElementById("setting_distance_area");
const distance_button = document.getElementById("setting_notify_distance");

distance_button.addEventListener("click", async function (evt) {
    //現在の通知距離取得
    const distance = await GetDistance()
    document.getElementById("distance_show").innerText = `現在は${distance}mで通知します`;
   
    // ポップアップ表示
    popup_distance_area.classList.add("is-show");
});


document.getElementById("notify_distances_select").addEventListener("change", async function() {
    // JWT取得
    const jwtToken = await GetJwt();

    // 選択された距離を取得
    const selectedDistance = this.value;

    // POSTリクエストのボディを作成
    const bodyData = {
        "Distance": parseInt(selectedDistance)
    };

    // fetchを使用してPOSTリクエストを送信
    const req = await fetch('/app/notify/distance', {
        method: 'POST',
        headers: {
            "Authorization": jwtToken,
            "Content-type": "application/json"
        },
        body: JSON.stringify(bodyData),
    })

    //レスポンスがokの時、現在の距離を更新
    if (req.ok) {
        const distance = await GetDistance()
        document.getElementById("distance_show").innerText = `現在は${distance}mで通知します`;
    }

});


async function GetDistance() {
    // JWT取得
    const jwtToken = await GetJwt();

    const req = await fetch('/app/notify/distance', {
        method: 'GET',
        headers: {
            "Authorization": jwtToken,
            "Content-type": "application/json"
        },
    })

    //レスポンスがokの時、現在の距離を更新
    if(req.ok) {
        const result =  await req.json()
        return result["Distance"]
    }

    return -1
}