```mermaid
graph TD
    Merchant[ร้านค้า] -->|1. ยื่นกู้ผ่าน API /apply| API[API Gateway]
    Manager[Mananger] -->|4. กดอนุมัติ /approve| API
    
    subgraph Core System [ระบบแกนหลัก]
        API -->|2. Start Workflow| Temporal[Temporal Server]
        API -->|3. Save สัญญากู้| DB[(PostgreSQL)]
        API -->|4. ส่ง Signal บอก| Temporal
        
        Temporal <-->|5. เช็กยอดรายวัน & ตัดยอดออโต้| Worker[Go Loan Worker]
        Worker -->|อัปเดตยอดชำระ| DB
    end

    classDef APICSS fill:#baebe1,stroke:#007c88,stroke-width:2px,color:#000000;
    class API APICSS

    classDef DBCSS fill:#74c7eb,stroke:#0b00ff,stroke-width:2px,color:#000000;
    class DB DBCSS

    classDef WorkerCSS fill:#d9d48b,stroke:#c1b576,stroke-width:2px,color:#000000;
    class Worker WorkerCSS

    classDef Failure fill:#f8d7da,stroke:#dc3545,stroke-width:2px,color:#000000;
    class Temporal Failure;

    click Temporal href "javascript:void(0)" "6. ถ้า Workflow Failed จนครบ สามารถ Manual ด้วยระบบ Reset ได้"