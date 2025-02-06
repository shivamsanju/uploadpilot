# Determining Optimal Chunk Size for File Uploads

## Given Data:

- **Network Bandwidth**: 10 Gbps (10,000 Mbps)
- **Concurrent Users**: 100
- **EC2 RAM**: 16 GB

---

## Step 1: Network Bandwidth Allocation per User

Since 100 users are uploading simultaneously, we divide the total bandwidth by the number of users:  
10,000 Mbps / 100 = **100 Mbps per user**

---

## Step 2: Recommended Chunk Size Based on Network Speed

A good upload chunk size should complete in less than a second.

With 100 Mbps per user, the optimal chunk size is:  
100 Mbps \* 1 second = **12.5 MB**

A chunk size between **12 MB - 16 MB** is reasonable.

---

## Step 3: Memory Constraints

Each user holds a chunk in memory before uploading.

If each user has one chunk in memory at a time:  
16 MB \* 100 users = **1.6 GB** of memory usage.

Since the total RAM is **16 GB**, using **16 MB chunks** is manageable.

---

## **Final Recommendation**

- Use a **chunk size of 12-16 MB** per upload for best performance.
- If memory usage is high, reduce it to **8 MB**.

---

# Server Requirements for 500 Users with 5 MB Chunks

## **Given Data:**

- **Chunk size** = 10 MB
- **Concurrent users** = 500

---

## **Step 1: Network Bandwidth Calculation**

Each user uploads **10 MB** chunks, and we assume each upload takes 1 second.

Total bandwidth required:  
10 MB _ 500 users _ 8 (to convert MB to Mbps) = **40 Gbps**

### **Network Requirement:**

- **Server must support at least 20 Gbps** network bandwidth.

---

## **Step 2: Memory Requirement Calculation**

Each user holds **one chunk in memory** before uploading.

Total memory required:  
10 MB \* 500 users = **5 GB RAM**

### **Memory Requirement:**

- A minimum of **8 GB RAM** should be allocated for upload processing.
- **16 GB+ RAM** is recommended if the server also handles requests, database operations, and caching.

---

## **Final Server Requirements for 500 Users with 10MB Chunks**

| **Component**           | **Minimum Requirement**                                |
| ----------------------- | ------------------------------------------------------ |
| **Network Bandwidth**   | **40 Gbps**                                            |
| **Memory (RAM)**        | **At least 8 GB (recommended 16 GB+ for other tasks)** |
| **CPU**                 | **At least 4+ vCPUs (depends on processing needs)**    |
| **Storage (Disk IOPS)** | **SSD with high write speed for handling uploads**     |

---
