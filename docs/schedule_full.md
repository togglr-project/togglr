# Understanding Feature Enablement and Schedules

### How the master switch and schedules control a feature

Each feature has a **Master Enable** switch and (optionally) one or more **schedules** that automatically change the feature state.

**1. Master Enable (global switch)**

* **If Master Enable = OFF** → the feature is completely turned off. Schedules and rules are ignored.
* **If Master Enable = ON** → the feature is “live” and its state is determined either manually (if no schedules exist) or by schedules (if schedules are configured).

**2. Two schedule types**

* **Repeating schedule** (the “repeat” builder you see in the UI). Internally it is stored as a cron-like rule and a required duration. Exactly **one repeating schedule** is allowed per feature.
* **One-shot schedules** (fixed start/end intervals). You may create **one or several non-overlapping** one-shot schedules for a feature. The UI prevents overlaps.

**3. How the current state is determined (algorithm)**
When Master Enable = ON:

* If **no schedules** are configured → the feature remains in its manual state (what you set in the UI). Schedules do not change anything because none exist.
* If **a repeating schedule** exists (the cron-like mode) → the feature’s baseline and active windows depend on the repeating schedule action:

    * If the schedule action is **Activate** → the baseline (outside scheduled windows) is **OFF**; the feature becomes **ON** only during scheduled windows (each window lasts the configured duration).
    * If the schedule action is **Deactivate** → the baseline (outside scheduled windows) is **ON**; the feature becomes **OFF** only during scheduled windows.
* If **one-shot schedules** exist (one or more intervals) → the baseline is derived from the collection of intervals:

    * If **all** one-shot intervals are `Activate` → baseline = **OFF** (feature is OFF except during the activate intervals).
    * If **any** one-shot interval is `Deactivate` (or there is a mix of Activate/Deactivate) → baseline = **ON** (feature is ON except during the deactivate intervals).
    * At any moment if a one-shot interval is active, that interval’s action (Activate => ON, Deactivate => OFF) determines visibility for that moment.

**4. Conflicts and precedence**

* You cannot mix repeating (cron) and one-shot schedules for the same feature (the UI and DB prevent that).
* If multiple schedules somehow apply at the same time, the system picks the schedule with the **latest creation time** (`created_at`). If two schedules have the same creation time, **Deactivate** wins over **Activate**.

**5. Required fields and notes**

* Every repeating schedule must include a **duration** (how long the action lasts after each trigger).
* Timezone is required and determines the local time at which scheduled actions run.
* The UI shows a friendly preview (human description + timeline) so you can confirm what you created before saving.

---

## Examples

* Feature `X` with Master Enable = OFF → `X` is always off, regardless of schedules.
* Feature `Y` with Master Enable = ON and **no schedules** → `Y` stays in whatever manual state you set in the UI.
* Feature `Z` with Master Enable = ON and **one repeating schedule** action = Activate, duration = 30m (daily at 09:30) → `Z` is OFF by default and turns ON for 30 minutes starting at 09:30 local time.
* Feature `A` with Master Enable = ON and two one-shot intervals: (`Deactivate` on Sep 20 18:00–18:30) and (`Activate` on Sep 25 09:00–10:00) → baseline = ON (because a Deactivate interval exists). At Sep 20 18:05 → `A` is OFF. At Sep 25 09:15 → `A` is ON. At other times → `A` is ON.
